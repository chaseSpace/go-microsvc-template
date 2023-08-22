package svccli

import (
	"microsvc/consts"
	"microsvc/infra/svcdiscovery"
	"microsvc/protocol/svc/admin"
	"microsvc/protocol/svc/user"
	"sync"
)

//type IntCliAPI interface {
//	Stop()
//}
//var _ IntCliAPI = new(intCli)

type intCli struct {
	once sync.Once
	inst *svcdiscovery.InstanceImpl
	svc  consts.Svc

	userCli  user.UserIntClient
	adminCli admin.AdminSvcClient
}

func (i *intCli) Stop() {
	if i.inst != nil {
		i.inst.Stop()
	}
}

func newIntCli(svc consts.Svc) *intCli {
	cli := &intCli{svc: svc}
	initializedSvc = append(initializedSvc, cli)
	return cli
}

func User() user.UserIntClient {
	userSvc.once.Do(func() {
		userSvc.inst = svcdiscovery.NewInstance(consts.SvcUser.Name())
	})
	v, err := userSvc.inst.GetInstance()
	if err == nil {
		userSvc.userCli = user.NewUserIntClient(v.Conn)
	}
	return user.NewUserIntClient(newFailGrpcClientConn())
}

func Admin() admin.AdminSvcClient {
	userSvc.once.Do(func() {
		adminSvc.inst = svcdiscovery.NewInstance(consts.SvcAdmin.Name())
	})
	v, err := adminSvc.inst.GetInstance()
	if err == nil {
		adminSvc.adminCli = admin.NewAdminSvcClient(v.Conn)
	}
	return admin.NewAdminSvcClient(newFailGrpcClientConn())
}
