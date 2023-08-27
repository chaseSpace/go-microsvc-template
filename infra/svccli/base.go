package svccli

import (
	"context"
	"google.golang.org/grpc"
	"microsvc/consts"
	"microsvc/deploy"
	"microsvc/infra/svcdiscovery"
	"microsvc/infra/svcdiscovery/sd"
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

type intCli struct {
	once      sync.Once
	svc       consts.Svc
	inst      *sd.InstanceImpl
	genClient sd.GenClient
}

var emptyConn = newFailGrpcClientConn()

func newIntCli(svc consts.Svc, gc sd.GenClient) *intCli {
	cli := &intCli{svc: svc, genClient: gc}
	return cli
}

func (ic *intCli) Getter() any {
	ic.once.Do(func() {
		ic.inst = sd.NewInstance(ic.svc.Name(), ic.genClient, svcdiscovery.GetSD())
		initializedSvcCli = append(initializedSvcCli, ic)
	})
	v, err := ic.inst.GetInstance()
	if err == nil {
		return v.Client
	}
	return ic.genClient(emptyConn)
}

func (i *intCli) Stop() {
	if i.inst != nil {
		i.inst.Stop()
	}
}

var initializedSvcCli []*intCli

func Stop() {
	for _, svcCli := range initializedSvcCli {
		svcCli.Stop()
	}
	xlog.Debug("svccli: resource released...")
}

func newFailGrpcClientConn() *grpc.ClientConn {
	cc, _ := grpc.DialContext(context.Background(), "127.0.0.1:1")
	return cc
}
