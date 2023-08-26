package main

import (
	"google.golang.org/grpc"
	"microsvc/deploy"
	"microsvc/infra"
	"microsvc/infra/svccli"
	"microsvc/infra/svcdiscovery"
	"microsvc/infra/xgrpc"
	"microsvc/pkg"
	"microsvc/pkg/xlog"
	"microsvc/protocol/svc/user"
	deploy2 "microsvc/service/user/deploy"
	"microsvc/service/user/handler"
	"microsvc/util/graceful"
)

func main() {
	// 初始化config
	deploy.Init("user", deploy2.UserConf)
	// 初始化服务用到的基础组件（封装于pkg目录下），如log, kafka等
	pkg.Init(
		xlog.Init,
		// 假如我要新增kafka等组件，也是新增 pkg/xkafka目录，然后实现其init函数并添加在这里
	)

	// 初始化几乎每个服务都需要的infra组件，must参数指定是否必须初始化成功，若must=true且err非空则panic
	infra.MustSetup(
		graceful.Init(),
		//cache.InitRedis(true),
		//orm.InitGorm(true),
		svcdiscovery.Init(true, deploy2.UserConf),
		svccli.Init(true),
	)

	x := xgrpc.New() // New一个封装好的grpc对象
	x.Apply(func(s *grpc.Server) {
		// 注册外部和内部的rpc接口对象
		user.RegisterUserExtServer(s, new(handler.UserExtCtrl))
		user.RegisterUserIntServer(s, new(handler.UserIntCtrl))
	})

	x.SetHTTPRegister(user.RegisterUserExtHandler) // 为外部接口对象 UserExt 启用 http反向代理 （http --call--> grpc）

	go graceful.OnExit(func() {
		x.Stop()
		infra.Stop()
	})

	x.Serve() // 监听请求
}
