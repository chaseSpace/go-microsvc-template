package abstract

import (
	"container/list"
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"microsvc/pkg/xerr"
	"microsvc/pkg/xlog"
	"time"
)

const logPrefix = "sd: "

type InstanceImpl struct {
	svc        string
	entryCache map[string]*GrpcConnObj
	grpcConns  *list.List    // 链表
	curr       *list.Element // 当前元素
	quit       chan struct{}
	genClient  GenClient
	sd         ServiceDiscovery
}

type GrpcConnObj struct {
	addr   string
	cc     *grpc.ClientConn
	Client interface{}
}

type GenClient func(conn *grpc.ClientConn) interface{}

func NewInstance(svc string, genClient GenClient, discovery ServiceDiscovery) *InstanceImpl {
	ins := &InstanceImpl{
		svc:        svc,
		entryCache: make(map[string]*GrpcConnObj),
		grpcConns:  list.New(),
		genClient:  genClient,
		quit:       make(chan struct{}),
		sd:         discovery,
	}
	_ = ins.refresh()
	go ins.backgroundRefresh()
	return ins
}

// GetInstance 每次返回链表的下一个元素，实现负载均衡（conn）
func (i *InstanceImpl) GetInstance() (inst *GrpcConnObj, err error) {
	time.Sleep(time.Second)
	if i.curr != nil {
		obj := i.curr.Value.(*GrpcConnObj)
		i.curr = i.curr.Next()
		return obj, nil
	}
	if elem := i.grpcConns.Front(); elem != nil { // 第一个
		obj := elem.Value.(*GrpcConnObj)
		i.curr = elem.Next()
		return obj, nil
	}
	return nil, xerr.ErrInternal.NewMsg(logPrefix+"%s no instance available.", i.svc)
}

func (i *InstanceImpl) backgroundRefresh() {
	for {
		err := i.refresh()
		select {
		case <-i.quit:
			xlog.Debug(logPrefix+"quited", zap.String("Svc", i.svc))
			return
		default:
			if err != nil {
				xlog.Error(logPrefix+"refresh err, hold on...", zap.Error(err))
				time.Sleep(time.Second * 3)
			}
		}

	}
}

func (i *InstanceImpl) refresh() error {
	var (
		entries []ServiceInstance
		cc      *grpc.ClientConn
		err     error
		ctx     context.Context
	)
	discovery := func() ([]ServiceInstance, error) {
		ctx = context.WithValue(context.Background(), CtxDurKey{}, time.Minute*2)
		return i.sd.Discover(ctx, i.svc)
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
	var availableEntries = make(map[string]int8)
	for _, entry := range entries {
		addr := entry.Addr()
		// check if entry is exists
		if i.entryCache[addr] != nil {
			continue
		}
		cc, err = newGRPCClient(addr)
		if err == nil {
			xlog.Debug(logPrefix+"newGRPCClient", zap.String("addr", addr))

			//println(2222, cc, addr)
			obj := &GrpcConnObj{addr: addr, Client: i.genClient(cc), cc: cc}
			availableEntries[addr] = 1
			i.entryCache[addr] = obj
			i.grpcConns.PushBack(obj)
		} else {
			xlog.Error(logPrefix+"newGRPCClient failed", zap.Error(err))
		}
	}

	// clear unavailable entries
	for addr, conn := range i.entryCache {
		if availableEntries[addr] == 0 {
			_ = conn.cc.Close()
			delete(i.entryCache, addr)
			i.removeGRPCConn(addr)
			xlog.Debug(logPrefix+"removeGRPCConn", zap.String("addr", addr))
		}
	}
	return nil
}

func (i *InstanceImpl) Stop() {
	close(i.quit)
}

func (i *InstanceImpl) removeGRPCConn(addr string) {
	curr := i.grpcConns.Front()
	for curr != nil {
		if curr.Value.(*GrpcConnObj).addr == addr {
			i.grpcConns.Remove(curr)
			return
		}
		curr = curr.Next()
	}
}

func newGRPCClient(target string) (cc *grpc.ClientConn, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	cc, err = grpc.DialContext(ctx, target, grpc.WithInsecure())
	return
}
