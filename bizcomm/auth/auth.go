package auth

import (
	"microsvc/consts"
	"microsvc/enums"
	"time"
)

const (
	HeaderKey = "Authorization"
)

type CtxAuthenticated struct{}

type SvcUser struct {
	ExternalUID int64 `json:"ext-uid"`
	AuthenticatedUser
}

func (a *SvcUser) IsValid() bool {
	if a.AuthenticatedUser.IsValid() {
		return a.ExternalUID > 0
	}
	return false
}

type AdminUser struct {
	AuthenticatedUser
}

type AuthenticatedUser struct {
	Uid      int64           `json:"uid"`
	Sex      enums.Sex       `json:"sex"`
	LoginAt  consts.Datetime `json:"login-at"`
	RegAt    consts.Datetime `json:"reg-at"`
	_LoginAt time.Time
	_RegAt   time.Time
}

func (a *AuthenticatedUser) IsValid() bool {
	var err1, err2 error
	a._LoginAt, err1 = a.LoginAt.Time()
	a._RegAt, err1 = a.RegAt.Time()
	if err1 != nil || err2 != nil {
		return false
	}
	if a.Uid > 0 &&
		a.Sex > enums.SexUnknown {
		return true
	}
	return false
}

func (a *AuthenticatedUser) GetLoginAt() time.Time {
	return a._LoginAt
}

func (a *AuthenticatedUser) GetRegAt() time.Time {
	return a._RegAt
}
