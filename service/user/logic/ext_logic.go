package logic

import (
	"github.com/golang-jwt/jwt/v5"
	"microsvc/bizcomm/auth"
	"microsvc/consts"
	"microsvc/deploy"
	"microsvc/enums"
	"microsvc/infra/xgrpc"
	"time"
)

func GenLoginToken(uid, extUID int64, sex enums.Sex) (string, error) {
	now := time.Now()
	nowDatetime := consts.Datetime(now.Format(time.DateTime))

	token, err := auth.GenerateJwT(
		xgrpc.SvcClaims{
			Authenticated: auth.SvcUser{
				ExternalUID: extUID,
				AuthenticatedUser: auth.AuthenticatedUser{
					Uid:     uid,
					Sex:     sex,
					LoginAt: nowDatetime,
					RegAt:   nowDatetime,
				},
			},
			RegisteredClaims: jwt.RegisteredClaims{},
		}, deploy.XConf.SvcTokenSignKey)

	return token, err
}
