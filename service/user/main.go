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
	dao "microsvc/proto/model/user"
	"microsvc/protocol/svc/user"
	"microsvc/service/user/handler"
)

func main() {
	_ = deploy.XConf // 初始化config

	dao.MustReady() // 检查dao层是否准备就绪

	// 初始化各类infra组件，must参数指定是否必须初始化成功，若must=true且err非空则panic
	infra.MustSetup(
		cache.InitRedis(false),
		orm.InitGorm(true),
		svcregistar.Init(true),
		svccli.Init(true),
	)

	// 创建一个grpc svr，并配置适当的中间件
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(
		xgrpc.RecoverGrpcRequest,
		xgrpc.LogGrpcRequest,
	))

	// 注册外部和内部的rpc接口组
	user.RegisterUserExtServer(server, new(handler.UserExtCtrl))
	user.RegisterUserIntServer(server, new(handler.UserIntCtrl))

	// 启动grpc服务
	xgrpc.ServeGRPC(server)
}
