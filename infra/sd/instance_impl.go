package sd

import (
	"container/list"
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"microsvc/infra/sd/abstract"
	"microsvc/infra/xgrpc"
	"microsvc/pkg/xerr"
	"microsvc/pkg/xlog"
	"time"
)

type InstanceImpl struct {
	svc        string
	entryCache map[string]*GrpcInstance
	grpcConns  *list.List    // linked list
	curr       *list.Element // current element
	quit       chan struct{}
	genClient  GenClient
	sd         abstract.ServiceDiscovery
}

type GrpcInstance struct {
	addr      string
	Conn      *grpc.ClientConn
	RpcClient interface{}
}

type GenClient func(conn *grpc.ClientConn) interface{}

func NewInstance(svc string, genClient GenClient, discovery abstract.ServiceDiscovery) *InstanceImpl {
	ins := &InstanceImpl{
		svc:        svc,
		entryCache: make(map[string]*GrpcInstance),
		grpcConns:  list.New(),
		genClient:  genClient,
		quit:       make(chan struct{}),
		sd:         discovery,
	}
	//_ = ins.query() // activate
	go ins.backgroundRefresh()
	return ins
}

// GetInstance get next conn, here implement load balancing（svc node）
func (i *InstanceImpl) GetInstance() (instance *GrpcInstance, err error) {
	instance, err = i.getCurr()
	if err != nil || instance != nil {
		return
	}
	// linked list is empty, try to refresh without blocking
	_ = i.query(false)
	return i.getCurr()
}

func (i *InstanceImpl) getCurr() (instance *GrpcInstance, err error) {
	for i.curr != nil {
		instance = i.curr.Value.(*GrpcInstance)
		// then we move to next or first element
		if next := i.curr.Next(); next != nil {
			i.curr = next
		} else {
			i.curr = i.grpcConns.Front()
		}
		if i.isConnReady(instance) {
			return
		}
		instance = nil
	}
	err = xerr.ErrInternal.NewMsg(logPrefix+"%s no instance available", i.svc)
	return
}

func (i *InstanceImpl) isConnReady(instance *GrpcInstance) bool {
	// if conn is idle, connect it
	if instance.Conn.GetState() == connectivity.Idle {
		instance.Conn.Connect()
		return true
	}
	// if conn is shutdown(contains closing), remove it then try next one in outside.
	if instance.Conn.GetState() == connectivity.Shutdown {
		delete(i.entryCache, instance.addr)
		i.removeInstance(instance.addr)
		return false
	}
	return true
}

func (i *InstanceImpl) backgroundRefresh() {
	for {
		err := i.query(true)
		select {
		case <-i.quit:
			xlog.Debug(logPrefix+"quited", zap.String("Svc", i.svc))
			return
		default:
			if err != nil {
				xlog.Error(logPrefix+"query err, hold on...", zap.Error(err))
				time.Sleep(time.Second * 3)
			}
		}
	}
}

// 阻塞刷新（首次请求不阻塞）
func (i *InstanceImpl) query(block bool) error {
	var (
		entries []abstract.ServiceInstance
		cc      *grpc.ClientConn
		err     error
		ctx     context.Context
	)
	discovery := func() ([]abstract.ServiceInstance, error) {
		ctx = context.WithValue(context.Background(), abstract.CtxDurKey{}, time.Minute*2)
		return i.sd.Discovery(ctx, i.svc, block)
	}
	entries, err = discovery()
	if err != nil {
		if err == context.DeadlineExceeded {
			xlog.Debug(logPrefix + "discover timeout")
		} else {
			xlog.Error(logPrefix+"discover fail", zap.Error(err))
		}
		return err
	}
	xlog.Debug(logPrefix+"discover result", zap.Any("entries", entries))

	var availableEntries = make(map[string]int8)
	for _, entry := range entries {
		addr := entry.Addr()
		availableEntries[addr] = 1
		if obj := i.entryCache[addr]; obj != nil {
			continue
		}
		cc, err = xgrpc.NewGRPCClient(addr, i.svc)
		if err == nil {
			xlog.Debug(logPrefix+"newGRPCClient OK", zap.String("addr", addr))
			obj := &GrpcInstance{addr: addr, RpcClient: i.genClient(cc), Conn: cc}
			i.entryCache[addr] = obj
			i.grpcConns.PushBack(obj)
			if i.curr == nil {
				i.curr = i.grpcConns.Front()
			}
		} else {
			xlog.Error(logPrefix+"newGRPCClient failed", zap.Error(err))
		}
	}

	// clear unavailable entries
	for addr, conn := range i.entryCache {
		if availableEntries[addr] == 0 {
			_ = conn.Conn.Close()
			delete(i.entryCache, addr)
			i.removeInstance(addr)
			xlog.Debug(logPrefix+"removeInstance", zap.String("addr", addr))
		}
	}
	return nil
}

func (i *InstanceImpl) Stop() {
	close(i.quit)
}

func (i *InstanceImpl) removeInstance(addr string) {
	curr := i.grpcConns.Front()
	for curr != nil {
		if curr.Value.(*GrpcInstance).addr == addr {
			i.grpcConns.Remove(curr)
			return
		}
		curr = curr.Next()
	}
}
