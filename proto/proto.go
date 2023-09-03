package proto

import (
	"microsvc/pkg/xerr"
	"microsvc/protocol/svc"
)

type CommonRes struct {
	svc.BaseRes
	Data interface{} `json:"data"`
}

func RespondOK(data interface{}) *CommonRes {
	return &CommonRes{
		BaseRes: svc.BaseRes{
			Code: xerr.ErrNil.Code,
			Msg:  xerr.ErrNil.Msg,
		},
		Data: data,
	}
}
