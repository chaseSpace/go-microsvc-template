package handler

import (
	"context"
	"microsvc/protocol/svc/admin"
)

type AdminCtrl struct {
	admin.UnimplementedAdminSvcServer
}

var _ admin.AdminSvcServer = new(AdminCtrl)

func (a AdminCtrl) AdminLogin(ctx context.Context, req *admin.AdminLoginReq) (*admin.AdminLoginRsp, error) {
	return &admin.AdminLoginRsp{
		Token: "",
		UserInfo: &admin.LoginResBody{
			Uid:  0,
			Nick: "",
			Sex:  0,
		},
	}, nil
}
