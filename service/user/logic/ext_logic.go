package logic

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"microsvc/bizcomm/auth"
	"microsvc/deploy"
	"microsvc/enums"
	"microsvc/infra/xgrpc"
	"microsvc/pkg/xerr"
	"microsvc/pkg/xlog"
	"microsvc/proto/model/user"
	user2 "microsvc/protocol/svc/user"
	"microsvc/service/user/dao"
	"microsvc/util"
	"microsvc/util/db"
	"time"
)

func CreateUser(ctx context.Context, req *user2.SignUpReq) (userModel *user.User, err error) {
	userModel = new(user.User)

	birth, _ := time.ParseInLocation(time.DateOnly, req.Birthday, time.Local)

	userModel.Uid = 1 // 暂设为1，以便通过check
	userModel.Sex = enums.Sex(req.Sex)
	userModel.Nickname = req.Nickname
	userModel.Birthday = birth
	// TODO check phone verify code

	if err = userModel.Check(); err != nil {
		return nil, xerr.ErrParams.NewMsg(err.Error())
	}

	tryInsert := func(i int) (duplicate bool, err error) {
		// 搜索测试函数：TestConcurrencySignUp
		ctx, cancel := context.WithTimeout(ctx, time.Second) // 经测试，1s可以抗住约60~80个并发请求（连接本地mysql），基本足够使用，可根据实际情况调整
		defer cancel()

		_uid, err := uidGenerator.GenUid(ctx)
		if err != nil {
			if err == context.DeadlineExceeded {
				return true, nil
			}
			return false, err
		}
		uid := int64(_uid)
		userModel.Uid = uid

		xlog.Info(fmt.Sprintf("CreateUser trying no.%d", i), zap.Any("model", *userModel))

		// inserting
		err = dao.CreateUser(userModel)
		if db.IsMysqlDuplicateErr(err) {
			return true, nil
		}
		if err == nil {
			uidGenerator.UpdateStartUid(_uid)
		}
		return false, err
	}

	var duplicate bool
	var tried int

	duplicate, err = tryInsert(tried)
	if err != nil {
		return
	}

	if duplicate {
		err = xerr.New("太多人注册辣，隔几秒再试一下哦")

		xlog.Error("CreateUser failed finally", zap.Int("tried", tried), zap.Any("lastModel", *userModel))
		return
	}
	return
}

func GenLoginToken(uid int64, regTime time.Time, sex enums.Sex) (string, error) {
	now := time.Now()

	expiry := auth.GetTokenExpiry(deploy.XConf.Env.IsDev())
	var expiresAt *jwt.NumericDate
	if expiry > 0 {
		expiresAt = jwt.NewNumericDate(now.Add(expiry))
	}
	token, err := auth.GenerateJwT(
		xgrpc.SvcClaims{
			SvcUser: auth.SvcUser{
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
				Subject:   fmt.Sprintf("%d", uid),
				ID:        util.NewKsuid(),
			},
		}, deploy.XConf.SvcTokenSignKey)

	return token, err
}
