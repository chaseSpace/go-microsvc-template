package main

import (
	"google.golang.org/grpc"
	"microsvc/deploy"
	"microsvc/enums"
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

	deploy.Init(enums.SvcAdmin, deploy2.AdminConf)

	pkg.Init(
		xlog.Init,
	)

	infra.MustSetup(
		//cache.InitRedis(true),
		//orm.InitGorm(true),
		sd.Init(true),
		svccli.Init(true),
	)

	x := xgrpc.New()
	x.Apply(func(s *grpc.Server) {
		admin.RegisterAdminExtServer(s, new(handler.AdminCtrl))
	})

	// 仅开发环境需要启动HTTP端口来代理gRPC服务
	if deploy.XConf.IsDevEnv() {
		x.SetHTTPExtRegister(admin.RegisterAdminExtHandler)
	}

	x.Start(deploy.XConf)
	sd.Register(deploy.XConf)

	graceful.Run()
}
