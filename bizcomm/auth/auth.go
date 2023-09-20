package auth

import (
	"context"
	"microsvc/consts"
	"microsvc/enums"
	"time"
)

const (
	HeaderKey               = "Authorization"
	TokenExpiry             = time.Hour * 3
	ForeverValidTokenExpiry = 0
	TokenIssuer             = "x.microsvc"
)

func GetTokenExpiry(isDev bool) time.Duration {
	if isDev {
		return ForeverValidTokenExpiry
	}
	return TokenExpiry
}

type Method struct {
}

var NoAuthMethods = map[string]*Method{
	"/svc.user.UserExt/Signup": {},
}

type CtxAuthenticated struct{}

func GetAuthUser(ctx context.Context) SvcUser {
	return ctx.Value(CtxAuthenticated{}).(SvcUser)
}
func GetAuthAdminUser(ctx context.Context) AdminUser {
	return ctx.Value(CtxAuthenticated{}).(AdminUser)
}

type Authenticator interface {
	IsValid() bool
	GetUser() interface{}
}

type SvcUser struct {
	ExternalUID int64 `json:"ext-uid"`
	AuthenticatedUser
}

func (a SvcUser) IsValid() bool {
	if a.AuthenticatedUser.IsValid() {
		return a.ExternalUID > 0
	}
	return false
}
func (a SvcUser) GetUser() interface{} {
	return a
}

type AdminUser struct {
	AuthenticatedUser
}

func (a AdminUser) GetUser() interface{} {
	return a
}

type AuthenticatedUser struct {
	Uid      int64     `json:"uid"`
	Sex      enums.Sex `json:"sex"`
	LoginAt  string    `json:"login-at"`
	RegAt    string    `json:"reg-at"`
	_LoginAt time.Time
	_RegAt   time.Time
}

func (a *AuthenticatedUser) IsValid() bool {
	var err1, err2 error
	a._LoginAt, err1 = consts.Datetime(a.LoginAt).Time()
	a._RegAt, err1 = consts.Datetime(a.RegAt).Time()
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
