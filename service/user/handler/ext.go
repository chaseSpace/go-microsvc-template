package handler

import (
	"context"
	"microsvc/enums"
	"microsvc/pkg/xerr"
	muser "microsvc/proto/model/user"
	"microsvc/protocol/svc/user"
	"microsvc/service/user/logic"
)

type UserExtCtrl struct {
}

var _ user.UserExtServer = new(UserExtCtrl)

func (u UserExtCtrl) Signup(ctx context.Context, req *user.SignupReq) (*user.SignupRes, error) {
	userModel := muser.User{}

	sex := enums.Sex(req.Sex)
	if !sex.IsValid() {
		return nil, xerr.ErrParams.AppendMsg("sex")
	}
	//TODO gen uid extUID
	userModel.SetIntField(1, 1, sex)

	token, err := logic.GenLoginToken(1, 1, sex)
	if err != nil {
		return nil, err
	}
	res := &user.SignupRes{Token: token}
	return res, nil
}

func (u UserExtCtrl) Login(ctx context.Context, req *user.LoginReq) (*user.LoginRes, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserExtCtrl) GetUser(ctx context.Context, req *user.GetUserReq) (*user.GetUserRes, error) {
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
