package handler

import (
	"context"
	"microsvc/protocol/svc/user"
)

type UserIntCtrl struct {
}

var _ user.UserIntServer = new(UserIntCtrl)

func (u UserIntCtrl) Test(ctx context.Context, req *user.TestReq) (*user.TestRsp, error) {
	return &user.TestRsp{New: req.Old + 1}, nil
}

func (u UserIntCtrl) GetUser(ctx context.Context, req *user.GetUserIntReq) (*user.GetUserIntRsp, error) {
	//TODO implement me
	panic("implement me")
}
