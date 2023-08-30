package handler

import (
	"context"
	"encoding/json"
	"google.golang.org/protobuf/proto"
	"microsvc/pkg/xerr"
	"microsvc/protocol/svc"
	"microsvc/protocol/svc/user"
	"microsvc/util"
	"reflect"
)

type UserExtCtrl struct {
}

var _ user.UserExtServer = new(UserExtCtrl)

func (u UserExtCtrl) GatewayCall(ctx context.Context, req *svc.GatewayReq) (*svc.GatewayRsp, error) {
	op, ok := apiRegistration[req.ApiName]
	if !ok {
		return nil, xerr.ErrApiNotFound.AppendMsg(req.ApiName)
	}
	freq := op.Req()
	err := json.Unmarshal(req.Body, freq)
	if err != nil {
		return nil, xerr.ErrParams.AppendMsg("parse json failed: %v", err)
	}

	// hardcode has better performance than reflect
	var res proto.Message
	switch req.ApiName {
	case "GetUser":
		res, err = u.GetUser(ctx, freq.(*user.GetUserReq))
	}
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, xerr.ErrInternal.NewMsg("[%s] rpc reply is empty", req.ApiName)
	}
	grsp := &svc.GatewayRsp{Body: util.ToJson(res)}
	return grsp, nil
}

type operation struct {
	Handler reflect.Method
	Req     func() interface{}
}

var apiRegistration = map[string]operation{
	"GetUser": {Req: func() interface{} { return new(user.GetUserReq) }},
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
