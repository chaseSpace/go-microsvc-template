package svccli

import (
	"go.uber.org/zap"
	"microsvc/deploy"
	"microsvc/enums"
	"microsvc/infra/sd"
	"microsvc/infra/sd/abstract"
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
			defaultSD, err = sd.NewSD()
			if err != nil {
				xlog.Error("svccli: NewSD failed", zap.Error(err))
			}
		}
		onEnd(must, err)
	}
}

type rpcClient struct {
	once      sync.Once
	svc       enums.Svc
	inst      *abstract.InstanceImpl
	genClient abstract.GenClient
}

func NewCli(svc enums.Svc, gc abstract.GenClient) *rpcClient {
	cli := &rpcClient{svc: svc, genClient: gc}
	return cli
}

func (c *rpcClient) Getter() any {
	c.once.Do(func() {
		c.inst = abstract.NewInstance(c.svc.Name(), c.genClient, defaultSD)
		initializedSvcCli = append(initializedSvcCli, c)
	})
	v, err := c.inst.GetInstance()
	if err == nil {
		return v.Client
	}
	return c.genClient(xgrpc.NewInvalidGRPCConn(c.svc.Name()))
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
