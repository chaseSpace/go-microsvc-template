package handler

import (
	"context"
	"encoding/json"
	"google.golang.org/protobuf/proto"
	"microsvc/pkg/xerr"
	"microsvc/proto/api"
	"microsvc/protocol/svc"
	"microsvc/protocol/svc/user"
	"microsvc/util"
)

type UserExtCtrl struct {
}

var _ user.UserExtServer = new(UserExtCtrl)

func (u UserExtCtrl) GatewayCall(ctx context.Context, req *svc.ForwardReq) (grsp *svc.ForwardRes, err error) {
	op := api.LoadUserApi(req.Method)
	if op == nil {
		return nil, xerr.ErrApiNotFound.AppendMsg(req.Method)
	}
	freq := op.Req()
	err = json.Unmarshal(req.Body, freq)
	if err != nil {
		return nil, xerr.ErrParams.AppendMsg("parse json failed: %v", err)
	}
	var res proto.Message
	defer func() {
		if err != nil {
			return
		}
		if res == nil {
			err = xerr.ErrInternal.NewMsg("[%s] rpc reply is empty", req.Method)
		}
		grsp = &svc.ForwardRes{Body: util.ToJson(res)}
	}()

	// hardcode has better performance than reflect
	switch req.Method {
	case "GetUser":
		res, err = u.GetUser(ctx, freq.(*user.GetUserReq))
	default:
		err = xerr.ErrInternal.NewMsg("add api(%s) firstly", req.Method)
	}
	return
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
	rsp := &user.GetUserRes{Umap: map[int64]*user.User{
		1: &user.User{
			Uid:  1,
			Nick: "nic",
			Age:  3,
			Sex:  4,
		},
	}}
	return rsp, nil
}
