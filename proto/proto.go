package proto

import (
	"microsvc/pkg/xerr"
	"microsvc/protocol/svc"
)

type CommonRes struct {
	*svc.BaseRes
	Data        interface{} `json:"data,omitempty"`
	FromGateway bool        `json:"from_gateway,omitempty"`
}

// WrapExtResponse
/*
old external grpc response(err==nil):

	{
	  "a": 1
	}

new external grpc response:

	{
	   "code": 200,
	   "msg": "OK",
	   "data": {"a": 1},
	}
*/
func WrapExtResponse(data interface{}, err error, fromGateway bool) *CommonRes {
	base := &svc.BaseRes{
		Code: xerr.ErrNil.Code,
		Msg:  xerr.ErrNil.Msg,
	}
	if err != nil {
		xe := xerr.ToXErr(err)
		base.Code = xe.Code
		base.Msg = xe.Msg
	}
	return &CommonRes{
		BaseRes:     base,
		Data:        data,
		FromGateway: fromGateway,
	}
}
