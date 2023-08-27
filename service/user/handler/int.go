package handler

import (
	"context"
	"errors"
	"microsvc/protocol/svc"
	"microsvc/protocol/svc/user"
)

type UserIntCtrl struct {
}

var _ user.UserIntServer = new(UserIntCtrl)

func (u UserIntCtrl) Test(ctx context.Context, req *user.TestReq) (*user.TestRsp, error) {
	return &user.TestRsp{New: req.Old + 1}, nil
}

func (u UserIntCtrl) GetUser(ctx context.Context, req *user.GetUserIntReq) (*user.GetUserIntRsp, error) {
	if len(req.Uid) == 0 {
		return nil, errors.New("参数无效")
	}
	umap := make(map[int64]*user.IntUser)
	umap[1] = &user.IntUser{
		Uid:  1,
		Nick: "Luyi",
		Age:  18,
		Sex:  svc.Sex_Male,
	}
	return &user.GetUserIntRsp{Umap: umap}, nil
}
