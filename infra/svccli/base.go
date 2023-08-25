package svccli

import (
	"context"
	"google.golang.org/grpc"
	"microsvc/consts"
	"microsvc/deploy"
	"microsvc/infra/svcdiscovery"
)

func Init(must bool) func(*deploy.XConfig, func(must bool, err error)) {
	return func(cc *deploy.XConfig, onEnd func(must bool, err error)) {
		// todo
		var err error
		onEnd(must, err)
	}
}

type intCli struct {
	inst      *svcdiscovery.InstanceImpl
	genClient svcdiscovery.GenClient
}

var emptyConn = newFailGrpcClientConn()

func newIntCli(svc consts.Svc, gc svcdiscovery.GenClient) *intCli {
	cli := &intCli{inst: svcdiscovery.NewInstance(svc.Name(), gc), genClient: gc}
	initializedSvc = append(initializedSvc, cli)
	return cli
}

func (ic *intCli) Getter() any {
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

var initializedSvc []*intCli

func Stop() {
	for _, svcCli := range initializedSvc {
		svcCli.Stop()
	}
}

func newFailGrpcClientConn() *grpc.ClientConn {
	cc, _ := grpc.DialContext(context.Background(), "127.0.0.1:1")
	return cc
}
