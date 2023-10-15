package handler

import (
	"context"
	"fmt"
	"github.com/k0kubun/pp"
	"microsvc/bizcomm/auth"
	"microsvc/bizcomm/commuser"
	"microsvc/enums"
	"microsvc/pkg/xerr"
	"microsvc/protocol/svc/user"
	"microsvc/service/user/cache"
	"microsvc/service/user/dao"
	"microsvc/service/user/logic"
)

type UserExtCtrl struct {
}

var _ user.UserExtServer = new(UserExtCtrl)

func (UserExtCtrl) SignUp(ctx context.Context, req *user.SignUpReq) (*user.SignUpRes, error) {
	umodel, err := logic.CheckSignUpReq(req)
	if err != nil {
		return nil, err
	}
	err = logic.CreateUser(ctx, umodel)
	if err != nil {
		return nil, err
	}

	token, err := logic.GenLoginToken(umodel.Uid, umodel.CreatedAt, enums.Sex(req.Sex))
	if err != nil {
		return nil, err
	}
	res := &user.SignUpRes{Token: token}
	return res, nil
}

func (UserExtCtrl) SignIn(ctx context.Context, req *user.SignInReq) (*user.SignInRes, error) {
	err := logic.CheckSignInReq(req)
	if err != nil {
		return nil, err
	}
	xphone := commuser.GetDBPhone(req.PhoneAreaCode, req.Phone)
	_, umodel, err := dao.GetUserByPhone(xphone)
	if err != nil {
		return nil, err
	}
	fmt.Printf("111 phone: %v,   model:%+v\n", req.Phone, umodel)
	if umodel.Uid == 0 {
		return nil, xerr.ErrParams.New("手机号未注册")
	}
	token, err := logic.GenLoginToken(umodel.Uid, umodel.CreatedAt, umodel.Sex)
	if err != nil {
		return nil, err
	}
	res := &user.SignInRes{Token: token, Info: umodel.ToPb()}
	return res, nil
}

func (UserExtCtrl) GetUser(ctx context.Context, req *user.GetUserReq) (*user.GetUserRes, error) {
	// 可直接通过ctx获取已经认证的自己的基础信息
	u := auth.GetAuthUser(ctx)
	_, _ = pp.Println("GetUser self:", u)

	if len(req.Uids) == 0 {
		return nil, xerr.ErrParams.New("missing arg:`uids`")
	}
	umap, err := cache.GetUser(req.Uids...)
	if err != nil {
		return nil, err
	}
	rsp := &user.GetUserRes{Umap: make(map[int64]*user.User)}
	for _, i := range umap {
		rsp.Umap[i.Uid] = i.ToPb()
	}
	return rsp, nil
}
