package svccli

import (
	"go.uber.org/zap"
	"microsvc/deploy"
	"microsvc/enums"
	"microsvc/infra/sd"
	"microsvc/infra/sd/abstract"
	"microsvc/infra/sd/mdns"
	"microsvc/infra/xgrpc"
	"microsvc/pkg/xlog"
	"sync"
)

var defaultSD abstract.ServiceDiscovery

func SetDefaultSD(sd abstract.ServiceDiscovery) {
	defaultSD = sd
}

func Init(must bool) func(*deploy.XConfig, func(must bool, err error)) {
	return func(cc *deploy.XConfig, onEnd func(must bool, err error)) {
		var err error
		if defaultSD == nil {
			if cc.IsDevEnv() {
				defaultSD = mdns.New()
			} else {
				defaultSD, err = sd.NewConsulSD()
				if err != nil {
					xlog.Error("svccli: NewConsulSD failed", zap.Error(err))
				}
			}
		}
		onEnd(must, err)
	}
}

type RpcClient struct {
	once      sync.Once
	svc       enums.Svc
	inst      *sd.InstanceImpl
	genClient sd.GenClient
}

func NewCli(svc enums.Svc, gc sd.GenClient) *RpcClient {
	cli := &RpcClient{svc: svc, genClient: gc}
	return cli
}

// Getter returns gRPC Server Client
func (c *RpcClient) Getter() any {
	c.once.Do(func() {
		c.inst = sd.NewInstance(c.svc.Name(), c.genClient, defaultSD)
		initializedSvcCli = append(initializedSvcCli, c)
	})
	v, err := c.inst.GetInstance()
	if err == nil {
		return v.Client
	}
	return c.genClient(xgrpc.NewInvalidGRPCConn(c.svc.Name()))
}

func (c *RpcClient) Stop() {
	if c.inst != nil {
		c.inst.Stop()
	}
}

var initializedSvcCli []*RpcClient

func Stop() {
	for _, svcCli := range initializedSvcCli {
		svcCli.Stop()
	}
	xlog.Debug("svccli: resource released...")
}
