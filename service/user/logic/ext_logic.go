package logic

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"microsvc/bizcomm/auth"
	"microsvc/deploy"
	"microsvc/enums"
	"microsvc/infra/xgrpc"
	"microsvc/util"
	"time"
)

func GenLoginToken(uid, extUID int64, regTime time.Time, sex enums.Sex) (string, error) {
	now := time.Now()

	expiry := auth.GetTokenExpiry(deploy.XConf.Env.IsDev())
	var expiresAt *jwt.NumericDate
	if expiry > 0 {
		expiresAt = jwt.NewNumericDate(now.Add(expiry))
	}
	token, err := auth.GenerateJwT(
		xgrpc.SvcClaims{
			SvcUser: auth.SvcUser{
				ExternalUID: extUID,
				AuthenticatedUser: auth.AuthenticatedUser{
					Uid:     uid,
					Sex:     sex,
					LoginAt: now.Format(time.DateTime),
					RegAt:   regTime.Format(time.DateTime),
				},
			},
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: expiresAt,
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
				Issuer:    auth.TokenIssuer,
				Subject:   fmt.Sprintf("%d", extUID),
				ID:        util.NewKsuid(),
			},
		}, deploy.XConf.SvcTokenSignKey)

	return token, err
}
