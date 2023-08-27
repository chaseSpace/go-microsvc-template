package main

import (
	"google.golang.org/grpc"
	"microsvc/deploy"
	"microsvc/infra"
	"microsvc/infra/sd"
	"microsvc/infra/svccli"
	"microsvc/infra/xgrpc"
	"microsvc/pkg"
	"microsvc/pkg/xlog"
	"microsvc/protocol/svc/admin"
	deploy2 "microsvc/service/admin/deploy"
	"microsvc/service/admin/handler"
	"microsvc/util/graceful"
)

func main() {
	graceful.SetupSignal()
	defer graceful.OnExit()

	deploy.Init("admin", deploy2.AdminConf)

	pkg.Init(
		xlog.Init,
	)

	infra.MustSetup(
		//cache.InitRedis(true),
		//orm.InitGorm(true),
		sd.Init(true),
		svccli.Init(true),
	)

	x := xgrpc.New(xgrpc.WrapAdminRsp)
	x.Apply(func(s *grpc.Server) {
		admin.RegisterAdminSvcServer(s, new(handler.AdminCtrl))
	})

	// 仅开发环境需要启动HTTP端口来代理gRPC服务
	if deploy.XConf.IsDevEnv() {
		x.SetHTTPExtRegister(admin.RegisterAdminSvcHandler)
	}

	x.Start(deploy.XConf)
	sd.Register(deploy.XConf)

	graceful.Run()
}
