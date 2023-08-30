package abstract

import (
	"container/list"
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"microsvc/infra/xgrpc"
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
	Conn   *grpc.ClientConn
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
	_ = ins.blockRefresh() // activate
	go ins.backgroundRefresh()
	return ins
}

// GetInstance 每次返回链表的下一个元素，实现负载均衡（conn）
func (i *InstanceImpl) GetInstance() (obj *GrpcConnObj, err error) {
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
	return nil, xerr.ErrInternal.NewMsg(logPrefix+"%s no instance available", i.svc)
}

func (i *InstanceImpl) backgroundRefresh() {
	for {
		err := i.blockRefresh()
		select {
		case <-i.quit:
			xlog.Debug(logPrefix+"quited", zap.String("Svc", i.svc))
			return
		default:
			if err != nil {
				xlog.Error(logPrefix+"blockRefresh err, hold on...", zap.Error(err))
				time.Sleep(time.Second * 3)
			}
		}
	}
}

// 阻塞刷新（首次请求不阻塞）
func (i *InstanceImpl) blockRefresh() error {
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
			//println(2222, Conn, addr)
			obj := &GrpcConnObj{addr: addr, Client: i.genClient(cc), Conn: cc}
			i.entryCache[addr] = obj
			i.grpcConns.PushBack(obj)
		} else {
			xlog.Error(logPrefix+"newGRPCClient failed", zap.Error(err))
		}
	}

	// clear unavailable entries
	for addr, conn := range i.entryCache {
		if availableEntries[addr] == 0 {
			_ = conn.Conn.Close()
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
