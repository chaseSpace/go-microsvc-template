package handler

import (
	"context"
	"microsvc/protocol/svc/user"
)

type UserIntCtrl struct {
}

var _ user.UserIntServer = new(UserIntCtrl)

func (u UserIntCtrl) GetUser(ctx context.Context, req *user.GetUserIntReq) (*user.GetUserIntRsp, error) {
	//TODO implement me
	panic("implement me")
}
