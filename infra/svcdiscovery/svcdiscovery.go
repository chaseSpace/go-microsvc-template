package svcdiscovery

import (
	"container/list"
	"context"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"microsvc/deploy"
	"microsvc/infra/svcdiscovery/consul"
	"microsvc/infra/svcdiscovery/define"
	"microsvc/pkg/xerr"
	"microsvc/pkg/xlog"
	"time"
)

func Init(must bool) func(*deploy.XConfig, func(must bool, err error)) {
	return func(cc *deploy.XConfig, onEnd func(must bool, err error)) {
		// 在这里 决定使用 etcd/consul
		cli, err := consul.NewConsulSD()
		if err == nil {
			sd = cli
		} else {
			err = errors.Wrap(err, "NewConsulSD")
		}
		onEnd(must, err)
	}
}

const logPrefix = "svcdiscovery: "

// sd init at Init
var sd define.ServiceDiscovery

type InstanceImpl struct {
	svc        string
	entryCache map[string]*GrpcConnObj
	grpcConns  *list.List    // 链表
	curr       *list.Element // 当前元素
	quit       chan struct{}
	genClient  GenClient
}

type GrpcConnObj struct {
	ID     string
	cc     *grpc.ClientConn
	Client interface{}
}

type GenClient func(conn *grpc.ClientConn) interface{}

func NewInstance(svc string, genClient GenClient) *InstanceImpl {
	ins := &InstanceImpl{
		svc:       svc,
		grpcConns: list.New(),
		genClient: genClient,
		quit:      make(chan struct{}),
	}
	go ins.backgroundRefresh()
	return ins
}

// GetInstance 每次返回链表的下一个元素（conn）
func (i *InstanceImpl) GetInstance() (inst *GrpcConnObj, err error) {
	if i.curr != nil {
		obj := i.curr.Value.(*GrpcConnObj)
		i.curr = i.curr.Next()
		return obj, nil
	}
	if elem := i.grpcConns.Front(); elem != nil {
		obj := elem.Value.(*GrpcConnObj)
		i.curr = elem.Next()
		return obj, nil
	}
	return nil, xerr.ErrInternal.NewMsg(logPrefix+"%s no instance available.", i.svc)
}

func (i *InstanceImpl) backgroundRefresh() {
	var (
		entries []define.ServiceInstance
		cc      *grpc.ClientConn
		err     error
	)
	ctx := context.Background()

	discovery := func() ([]define.ServiceInstance, error) {
		ctx, cancel := context.WithTimeout(ctx, time.Minute)
		defer cancel()
		entries, err = sd.Discover(ctx, i.svc)
		return nil, err
	}
	for {
		entries, err = discovery()
		select {
		case <-i.quit:
			xlog.Debug(logPrefix+"quited", zap.String("svc", i.svc))
			return
		default:
		}
		if err != nil {
			if err == context.DeadlineExceeded {
				xlog.Debug(logPrefix + "Discover timeout")
			} else {
				xlog.Error(logPrefix+"Discover fail", zap.Error(err))
			}
			time.Sleep(time.Second * 3)
			continue
		}
		var availableEntries = make(map[string]int8)
		for _, entry := range entries {
			availableEntries[entry.ID] = 1
			// check if entry is exists
			if i.entryCache[entry.ID] != nil {
				continue
			}
			cc, err = newGrpcClient(entry.Addr())
			if err == nil {
				obj := &GrpcConnObj{ID: entry.Addr(), Client: i.genClient(cc), cc: cc}
				i.entryCache[entry.ID] = obj
				i.grpcConns.PushBack(obj)
			} else {
				xlog.Error(logPrefix+"newGrpcClient", zap.Error(err))
			}
		}

		// clear unavailable entries
		for svcId, conn := range i.entryCache {
			if availableEntries[svcId] == 0 {
				_ = conn.cc.Close()
				delete(i.entryCache, svcId)
				i.removeGrpcConn(svcId)
				xlog.Debug(logPrefix+"removeGrpcConn", zap.String("svcId", svcId))
			}
		}
	}
}

func (i *InstanceImpl) Stop() {
	close(i.quit)
}
func (i *InstanceImpl) removeGrpcConn(id string) {
	curr := i.grpcConns.Front()
	for curr != nil {
		if curr.Value.(*GrpcConnObj).ID == id {
			i.grpcConns.Remove(curr)
			return
		}
		curr = curr.Next()
	}
}

func newGrpcClient(target string) (cc *grpc.ClientConn, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	cc, err = grpc.DialContext(ctx, target, grpc.WithInsecure())
	return
}
