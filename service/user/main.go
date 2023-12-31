package main

import (
	"google.golang.org/grpc"
	"microsvc/deploy"
	"microsvc/enums/svc"
	"microsvc/infra"
	"microsvc/infra/sd"
	"microsvc/infra/svccli"
	"microsvc/infra/xgrpc"
	_ "microsvc/infra/xgrpc/protobytes"
	"microsvc/pkg"
	"microsvc/pkg/xkafka"
	"microsvc/pkg/xlog"
	"microsvc/protocol/svc/user"
	deploy2 "microsvc/service/user/deploy"
	"microsvc/service/user/handler"
	"microsvc/service/user/logic"
	"microsvc/util/graceful"
)

func main() {
	graceful.SetupSignal()
	defer graceful.OnExit()

	// 初始化config
	deploy.Init(svc.User, deploy2.UserConf)

	// 初始化服务用到的基础组件（封装于pkg目录下），如log, kafka等
	pkg.Setup(
		xlog.Init,
		xkafka.Init,
	)

	// 初始化几乎每个服务都需要的infra组件，must参数指定是否必须初始化成功，若must=true且err非空则panic
	// - 注意顺序
	infra.Setup(
		//cache.InitRedis(true),
		//orm.InitGorm(true),
		sd.Init(true),
		svccli.Init(true),
	)

	// 此服务需要的初始化(在infra初始化之后进行)
	logic.MustInit()

	x := xgrpc.New() // New一个封装好的grpc对象
	x.Apply(func(s *grpc.Server) {
		// 注册外部和内部的rpc接口对象
		user.RegisterUserExtServer(s, new(handler.UserExtCtrl))
		user.RegisterUserIntServer(s, new(handler.UserIntCtrl))
	})

	x.Start(deploy.XConf)

	sd.MustRegister(deploy.XConf)

	graceful.Run()
}
