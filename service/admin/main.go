package main

import (
	"google.golang.org/grpc"
	"microsvc/deploy"
	"microsvc/enums"
	"microsvc/infra"
	"microsvc/infra/sd"
	"microsvc/infra/svccli"
	"microsvc/infra/xgrpc"
	_ "microsvc/infra/xgrpc/protobytes"
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

	pkg.Setup(
		xlog.Init,
	)

	infra.Setup(
		//cache.InitRedis(true),
		//orm.InitGorm(true),
		sd.Init(true),
		svccli.Init(true),
	)

	x := xgrpc.New()
	x.Apply(func(s *grpc.Server) {
		admin.RegisterAdminExtServer(s, new(handler.AdminCtrl))
	})

	x.Start(deploy.XConf)
	sd.MustRegister(deploy.XConf)

	graceful.Run()
}
