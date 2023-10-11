package handler

import (
	"context"
	"microsvc/enums"
	"microsvc/pkg/xerr"
	"microsvc/protocol/svc/user"
	"microsvc/service/user/logic"
)

type UserExtCtrl struct {
}

var _ user.UserExtServer = new(UserExtCtrl)

func (UserExtCtrl) SignUp(ctx context.Context, req *user.SignUpReq) (*user.SignUpRes, error) {
	umodel, err := logic.CreateUser(ctx, req)
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
	//TODO implement me
	panic("implement me")
}

func (UserExtCtrl) GetUser(ctx context.Context, req *user.GetUserReq) (*user.GetUserRes, error) {
	//u := auth.GetAuthUser(ctx)
	//pp.Println(1111, u)
	if len(req.Uids) == 0 {
		return nil, xerr.ErrParams
	}
	//umap, err := cache.GetUser(req.Uids...)
	//if err != nil {
	//	return nil, err
	//}
	//rsp := &user.GetUserRes{Umap: make(map[int64]*user.User)}
	//for _, i := range umap {
	//	rsp.Umap[i.Uid] = i.ToPb()
	//}
	rsp := &user.GetUserRes{
		Umap: map[int64]*user.User{
			1: &user.User{
				Uid:      1,
				Nickname: "niko",
				Age:      3,
				Sex:      4,
			},
			2: &user.User{
				Uid:      2,
				Nickname: "lucy",
				Age:      3,
				Sex:      4,
			},
		}}
	return rsp, nil
}
