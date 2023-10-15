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
	"/svc.user.UserExt/SignUp": {},
	"/svc.user.UserExt/SignIn": {},
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
	AuthenticatedUser
}

func (a SvcUser) IsValid() bool {
	if a.AuthenticatedUser.IsValid() {
		return a.Uid > 0
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

func NewTestSvcUser(uid int64, sex enums.Sex) SvcUser {
	return SvcUser{
		AuthenticatedUser{
			Uid:     uid,
			Sex:     sex,
			RegAt:   "2024-01-01 00:00:00",
			LoginAt: time.Now().Format(time.DateTime),
		},
	}
}

func NewTestAdminUser(uid int64, sex enums.Sex) AdminUser {
	return AdminUser{
		AuthenticatedUser{
			Uid:     uid,
			Sex:     sex,
			RegAt:   "2024-01-01 00:00:00",
			LoginAt: time.Now().Format(time.DateTime),
		},
	}
}
