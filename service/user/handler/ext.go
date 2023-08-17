package handler

import (
	"context"
	"gorm.io/gorm"
	"microsvc/protocol/svc/user"
	"microsvc/service/user/dao"
)

type UserExtCtrl struct {
}

var _ user.UserExtServer = new(UserExtCtrl)

func (u UserExtCtrl) GetUser(ctx context.Context, req *user.GetUserReq) (*user.GetUserRsp, error) {
	list, _, err := dao.GetUser(req.Uids...)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &user.GetUserRsp{Umap: make(map[int64]*user.User)}, nil
		}
		return nil, err
	}
	rsp := &user.GetUserRsp{Umap: make(map[int64]*user.User)}
	for _, i := range list {
		rsp.Umap[i.Uid] = i.ToPb()
	}
	return rsp, nil
}
