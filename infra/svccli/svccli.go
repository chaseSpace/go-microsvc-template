package svccli

import (
	"google.golang.org/grpc"
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
		userSvc.inst = svcdiscovery.NewInstance(consts.SvcUser.Name(), func(conn *grpc.ClientConn) interface{} {
			return user.NewUserIntClient(conn)
		})
	})
	v, err := userSvc.inst.GetInstance()
	if err == nil {
		userSvc.userCli = v.Client.(user.UserIntClient)
	}
	return user.NewUserIntClient(newFailGrpcClientConn())
}

func Admin() admin.AdminSvcClient {
	userSvc.once.Do(func() {
		adminSvc.inst = svcdiscovery.NewInstance(consts.SvcAdmin.Name(), func(conn *grpc.ClientConn) interface{} {
			return admin.NewAdminSvcClient(conn)
		})
	})
	v, err := adminSvc.inst.GetInstance()
	if err == nil {
		adminSvc.adminCli = v.Client.(admin.AdminSvcClient)
	}
	return admin.NewAdminSvcClient(newFailGrpcClientConn())
}
