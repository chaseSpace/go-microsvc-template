package handler

import (
	"context"
	"go.uber.org/zap"
	"microsvc/infra/svccli"
	"microsvc/pkg/xlog"
	"microsvc/protocol/svc"
	"microsvc/protocol/svc/admin"
	"microsvc/protocol/svc/user"
)

type AdminCtrl struct {
}

var _ admin.AdminExtServer = new(AdminCtrl)

func (a AdminCtrl) AdminLogin(ctx context.Context, req *admin.AdminLoginReq) (*admin.AdminLoginRsp, error) {
	return &admin.AdminLoginRsp{
		Token: "token",
		UserInfo: &admin.LoginResBody{
			Uid:  123,
			Nick: "Luyi",
			Sex:  svc.Sex_Male,
		},
	}, nil
}

func (a AdminCtrl) GetUser(ctx context.Context, req *admin.GetUserReq) (*admin.GetUserRsp, error) {
	rsp, err := svccli.User().GetUser(ctx, &user.GetUserIntReq{
		Uid: req.Uid,
	})
	if err != nil {
		return nil, err
	}
	xlog.Info("111...", zap.Any("Umap", rsp.Umap))
	return &admin.GetUserRsp{Umap: rsp.Umap}, nil
}
