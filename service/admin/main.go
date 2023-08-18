package main

import (
	"google.golang.org/grpc"
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
	infra.MustSetup(
		cache.InitRedis(true),
		orm.InitGorm(true),
		svcregistar.Init(true),
		svccli.Init(true),
	)

	x := xgrpc.New(xgrpc.WrapAdminRsp)

	x.Apply(func(s *grpc.Server) {
		admin.RegisterAdminSvcServer(s, new(handler.AdminCtrl))
	})

	x.SetHTTPRegister(admin.RegisterAdminSvcHandler)
	x.Serve()
}
