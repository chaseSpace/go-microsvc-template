package svccli

import (
	"context"
	"google.golang.org/grpc"
	"microsvc/consts"
	"microsvc/deploy"
)

func Init(must bool) func(*deploy.XConfig, func(must bool, err error)) {
	return func(cc *deploy.XConfig, onEnd func(must bool, err error)) {
		// todo
		var err error
		onEnd(must, err)
	}
}

// 这里只会初始化 *IntClient , 即内部接口对象
var (
	userSvc  = newIntCli(consts.SvcUser)
	adminSvc = newIntCli(consts.SvcAdmin)
)

var initializedSvc []*intCli

func Stop() {
	for _, svcCli := range initializedSvc {
		svcCli.Stop()
	}
}

func newFailGrpcClientConn() *grpc.ClientConn {
	cc, _ := grpc.DialContext(context.Background(), "127.0.0.1")
	return cc
}
