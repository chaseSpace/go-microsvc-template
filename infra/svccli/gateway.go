//go:build !k8s

package svccli

import (
	"google.golang.org/grpc"
	"microsvc/enums/svc"
	"microsvc/infra/sd"
	"microsvc/infra/xgrpc"
	"sync"
)

type InstanceMgr struct {
	cmap map[svc.Svc]*InstanceImplT
	mu   sync.RWMutex
}

type InstanceImplT struct {
	impl   *sd.InstanceImpl
	errCnt int32
	mu     sync.RWMutex
}

var defaultConnMgr = &InstanceMgr{
	cmap: map[svc.Svc]*InstanceImplT{},
}

const cleanSvcInstanceErrCntThreshold = 10

// GetConn TODO: optimize, dont use global lock here
func GetConn(svc svc.Svc) (conn *grpc.ClientConn) {
	defaultConnMgr.mu.RLock()
	inst := defaultConnMgr.cmap[svc]
	defaultConnMgr.mu.RUnlock()

	defer func() {
		if conn == nil {
			conn = xgrpc.NewInvalidGRPCConn(svc.Name())
		}
	}()
	if inst != nil {
		obj, err := inst.impl.GetInstance()

		// segment lock
		inst.mu.Lock()
		defer inst.mu.Unlock()
		if err == nil {
			inst.errCnt = 0
			return obj.Conn
		} else if inst.errCnt > cleanSvcInstanceErrCntThreshold {
			defaultConnMgr.mu.Lock()
			delete(defaultConnMgr.cmap, svc)
			defaultConnMgr.mu.Unlock()
		} else {
			inst.errCnt++
		}
		return
	}

	// Add conn instance
	defaultConnMgr.mu.Lock()
	impl := sd.NewInstance(svc.Name(), func(conn *grpc.ClientConn) interface{} {
		return nil
	}, defaultSD)

	defaultConnMgr.cmap[svc] = &InstanceImplT{
		impl: impl,
		mu:   sync.RWMutex{},
	}
	defaultConnMgr.mu.Unlock()

	obj, err := impl.GetInstance()
	if err != nil {
		return nil
	}
	return obj.Conn
}
