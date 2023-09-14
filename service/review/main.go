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
	"microsvc/pkg/xkafka"
	"microsvc/pkg/xlog"
	"microsvc/protocol/svc/review"
	deploy2 "microsvc/service/review/deploy"
	"microsvc/service/review/handler"
	"microsvc/util/graceful"
)

func main() {
	graceful.SetupSignal()
	defer graceful.OnExit()

	// 初始化config
	deploy.Init(enums.SvcReview, deploy2.ReviewConf)

	// 初始化服务用到的基础组件（封装于pkg目录下），如log, kafka等
	pkg.Setup(
		xlog.Init,
		xkafka.Init,
	)

	// 初始化几乎每个服务都需要的infra组件，must参数指定是否必须初始化成功，若must=true且err非空则panic
	infra.Setup(
		//cache.InitRedis(true),
		//orm.InitGorm(true),
		sd.Init(true),
		svccli.Init(true),
	)

	x := xgrpc.New() // New一个封装好的grpc对象
	x.Apply(func(s *grpc.Server) {
		// 注册内部的rpc接口对象（review服务一般不对外提供接口，所以模板只有int接口对象）
		review.RegisterUserIntServer(s, new(handler.ReviewIntCtrl))
	})

	x.Start(deploy.XConf)
	// GRPC服务启动后 再注册服务
	sd.MustRegister(deploy.XConf)

	graceful.Run()
}
