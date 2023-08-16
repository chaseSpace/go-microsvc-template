package main

import (
	"google.golang.org/grpc"
	"microsvc/deploy"
	"microsvc/infra"
	"microsvc/infra/cache"
	"microsvc/infra/orm"
	"microsvc/infra/svccli"
	"microsvc/infra/svcregistar"
	"microsvc/infra/xgrpc"
	"microsvc/protocol/svc/admin"
	"microsvc/service/admin/handler"
)

func main() {
	_ = deploy.XConf // init config

	infra.MustSetup(
		cache.InitRedis(true),
		orm.InitGorm(true),
		svcregistar.Init(true),
		svccli.Init(true),
	)

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(
		xgrpc.WrapAdminRsp,
	))
	admin.RegisterAdminSvcServer(server, new(handler.AdminCtrl))

	xgrpc.ServeGRPC(server)
}
