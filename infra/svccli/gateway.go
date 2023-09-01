package svccli

import (
	"google.golang.org/grpc"
	"microsvc/enums/svc"
	"microsvc/infra/sd/abstract"
	"sync"
)

type ConnMgr struct {
	cmap map[svc.Svc]*InstanceImplT
	mu   sync.RWMutex
}

type InstanceImplT struct {
	impl   *abstract.InstanceImpl
	errCnt int32
}

var defaultConnMgr = &ConnMgr{
	cmap: map[svc.Svc]*InstanceImplT{},
}

const cleanSvcInstanceErrCntThreshold = 50

func GetConn(svc svc.Svc) *grpc.ClientConn {
	defaultConnMgr.mu.RLock()
	inst := defaultConnMgr.cmap[svc]
	defaultConnMgr.mu.RUnlock()

	if inst == nil {
		impl := abstract.NewInstance(svc.Name(), func(conn *grpc.ClientConn) interface{} {
			return nil
		}, defaultSD)

		obj, err := impl.GetInstance()
		if err != nil {
			return nil
		}
		defaultConnMgr.mu.Lock()
		defaultConnMgr.cmap[svc] = &InstanceImplT{impl: impl}
		defaultConnMgr.mu.Unlock()
		return obj.Conn
	}

	obj, err := inst.impl.GetInstance()
	if err == nil {
		return obj.Conn
	} else if inst.errCnt > cleanSvcInstanceErrCntThreshold {
		defaultConnMgr.mu.Lock()
		delete(defaultConnMgr.cmap, svc)
		defaultConnMgr.mu.Unlock()
	} else {
		inst.errCnt++
	}
	return nil
}
