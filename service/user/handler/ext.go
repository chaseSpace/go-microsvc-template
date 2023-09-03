package handler

import (
	"context"
	"microsvc/pkg/xerr"
	"microsvc/protocol/svc/user"
)

type UserExtCtrl struct {
}

var _ user.UserExtServer = new(UserExtCtrl)

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
				Uid:  1,
				Nick: "niko",
				Age:  3,
				Sex:  4,
			},
			2: &user.User{
				Uid:  2,
				Nick: "lucy",
				Age:  3,
				Sex:  4,
			},
		}}
	return rsp, nil
}
