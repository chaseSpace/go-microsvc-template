package main

import (
	"google.golang.org/grpc"
	"microsvc/infra"
	"microsvc/infra/cache"
	"microsvc/infra/orm"
	"microsvc/infra/svccli"
	"microsvc/infra/svcregistar"
	"microsvc/infra/xgrpc"
	"microsvc/pkg/xlog"
	dao "microsvc/proto/model/user"
	"microsvc/protocol/svc/user"
	"microsvc/service/user/deploy"
	"microsvc/service/user/handler"
)

func main() {
	// 初始化config
	deploy.MustSetup(
		// 在这里传入pkg内需要初始化的函数
		xlog.Init,
		// 假如我要新增kafka等组件，也是新增 pkg/xkafka目录，然后实现其init函数并添加在这里
	)

	// 初始化各类infra组件，must参数指定是否必须初始化成功，若must=true且err非空则panic
	infra.MustSetup(
		cache.InitRedis(false),
		orm.InitGorm(true),
		svcregistar.Init(true),
		svccli.Init(true),
	)
	dao.MustReady() // 检查dao层是否准备就绪

	// 创建一个grpc svr，并配置适当的中间件
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(
		xgrpc.RecoverGrpcRequest,
		xgrpc.LogGrpcRequest,
	))

	// 注册外部和内部的rpc接口组
	user.RegisterUserExtServer(server, new(handler.UserExtCtrl))
	user.RegisterUserIntServer(server, new(handler.UserIntCtrl))

	// 启动grpc服务
	x := xgrpc.New(server)
	// -- 为UserExt启用 http反向代理 （http --call--> grpc）
	x.SetHTTPRegister(user.RegisterUserExtHandler)
	x.Serve()
}
