package svccli

import (
	"google.golang.org/grpc"
	"microsvc/deploy"
	"microsvc/enums"
	"microsvc/infra/sd"
	"microsvc/infra/sd/abstract"
	"microsvc/pkg/xlog"
	"sync"
)

func Init(must bool) func(*deploy.XConfig, func(must bool, err error)) {
	return func(cc *deploy.XConfig, onEnd func(must bool, err error)) {
		// todo
		var err error
		onEnd(must, err)
	}
}

type rpcClient struct {
	once      sync.Once
	svc       enums.Svc
	inst      *abstract.InstanceImpl
	genClient abstract.GenClient
}

var emptyConn = newFailGrpcClientConn()

func NewCli(svc enums.Svc, gc abstract.GenClient) *rpcClient {
	cli := &rpcClient{svc: svc, genClient: gc}
	return cli
}

func (c *rpcClient) Getter() any {
	c.once.Do(func() {
		c.inst = abstract.NewInstance(c.svc.Name(), c.genClient, sd.GetSD())
		initializedSvcCli = append(initializedSvcCli, c)
	})
	v, err := c.inst.GetInstance()
	if err == nil {
		return v.Client
	}
	return c.genClient(emptyConn)
}

func (c *rpcClient) Stop() {
	if c.inst != nil {
		c.inst.Stop()
	}
}

var initializedSvcCli []*rpcClient

func Stop() {
	for _, svcCli := range initializedSvcCli {
		svcCli.Stop()
	}
	xlog.Debug("svccli: resource released...")
}

func newFailGrpcClientConn() *grpc.ClientConn {
	cc := &grpc.ClientConn{}
	return cc
}
