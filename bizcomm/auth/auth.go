package auth

import (
	"microsvc/consts"
	"microsvc/enums"
)

const (
	HttpHeaderKey = "Authorization"
)

type CtxAuthenticated struct{}

type AuthenticatedSvcUser struct {
	ExternalUID int64           `json:"ext-uid"`
	InternalUID int64           `json:"int-uid"`
	Sex         enums.Sex       `json:"sex"`
	LoginAt     consts.Datetime `json:"login-at"`
	RegAt       consts.Datetime `json:"reg-at"`
}

type AuthenticatedAdminUser struct {
	Uid      int64           `json:"uid"`
	Nickname string          `json:"nickname"`
	Sex      enums.Sex       `json:"sex"`
	LoginAt  consts.Datetime `json:"login-at"`
	RegAt    consts.Datetime `json:"reg-at"`
}
