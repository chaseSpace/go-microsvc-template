package main

import (
	"google.golang.org/grpc"
	"microsvc/consts"
	"microsvc/deploy"
	"microsvc/infra"
	"microsvc/infra/cache"
	"microsvc/infra/orm"
	"microsvc/infra/svccli"
	"microsvc/infra/svcdiscovery"
	"microsvc/infra/xgrpc"
	"microsvc/protocol/svc/admin"
	deploy2 "microsvc/service/admin/deploy"
	"microsvc/service/admin/handler"
)

func main() {
	svc := consts.SvcAdmin
	deploy.Init(svc, deploy2.AdminConf)

	infra.MustSetup(
		cache.InitRedis(true),
		orm.InitGorm(true),
		svcdiscovery.Init(true),
		svccli.Init(true),
	)
	{
		svcdiscovery.GetSD().Register(svc.Name(), "", 1, nil)
	}
	defer func() {
		_ = svcdiscovery.GetSD().Deregister(svc.Name())
		infra.Stop()
	}()

	x := xgrpc.New(xgrpc.WrapAdminRsp)
	x.Apply(func(s *grpc.Server) {
		admin.RegisterAdminSvcServer(s, new(handler.AdminCtrl))
	})
	x.SetHTTPRegister(admin.RegisterAdminSvcHandler)
	x.Serve()
}
