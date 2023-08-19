package handler

import (
	"context"
	"microsvc/protocol/svc/user"
	"microsvc/service/user/cache"
)

type UserExtCtrl struct {
}

var _ user.UserExtServer = new(UserExtCtrl)

func (u UserExtCtrl) GetUser(ctx context.Context, req *user.GetUserReq) (*user.GetUserRsp, error) {
	umap, err := cache.GetUser(req.Uids...)
	if err != nil {
		return nil, err
	}
	rsp := &user.GetUserRsp{Umap: make(map[int64]*user.User)}
	for _, i := range umap {
		rsp.Umap[i.Uid] = i.ToPb()
	}
	return rsp, nil
}
