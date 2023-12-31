package logic

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"microsvc/bizcomm/auth"
	"microsvc/bizcomm/commuser"
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

func CreateUser(ctx context.Context, userModel *user.User) (err error) {
	tryInsert := func() (duplicate bool, err error) {
		// 搜索测试函数：TestConcurrencySignUp
		ctx, cancel := context.WithTimeout(ctx, time.Second)
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

		xlog.Info("CreateUser before", zap.Any("model", *userModel))

		// inserting
		err = dao.CreateUser(userModel)

		forIndex := ""
		if db.IsMysqlDuplicateErr(err, &forIndex) {
			if forIndex == user.UniqueKeyPhone {
				return false, xerr.ErrParams.New("该手机号已被注册")
			}
			if forIndex != user.UniqueKeyUID {
				return false, xerr.New("user表其他唯一列出现重复，请检查")
			}
			println(1212, uid)
			return true, nil
		}
		return false, err
	}

	var duplicate bool

	duplicate, err = tryInsert()
	if err != nil {
		return
	}

	if duplicate {
		err = xerr.New("太多人注册辣，隔几秒再试一下哦")

		xlog.Error("CreateUser failed finally", zap.Any("lastModel", *userModel))
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

func CheckSignUpReq(req *user2.SignUpReq) (*user.User, error) {
	birth, _ := time.ParseInLocation(time.DateOnly, req.Birthday, time.Local)

	umodel := &user.User{
		Base: user.Base{
			Uid:      1,
			Nickname: req.Nickname,
			Birthday: birth,
			Sex:      enums.Sex(req.Sex),
			Phone:    commuser.GetDBPhone(req.PhoneAreaCode, req.Phone),
		},
	}
	err := umodel.Check()
	if err != nil {
		return nil, xerr.ErrParams.New(err.Error())
	}

	if !commuser.IsPhoneAreaCodeSupported(req.PhoneAreaCode) {
		return nil, xerr.ErrParams.New("请提供有效的手机区号")
	}
	switch req.PhoneAreaCode {
	case "86":
		if len(req.Phone) != 11 {
			return nil, xerr.ErrParams.New("请提供有效的手机号")
		}
	}
	if len(req.VerifyCode) != 4 {
		return nil, xerr.ErrParams.New("请提供有效的验证码")
	}
	// TODO: check req.VerifyCode

	umodel.Uid = 0
	return umodel, nil
}

func CheckSignInReq(req *user2.SignInReq) error {
	if !commuser.IsPhoneAreaCodeSupported(req.PhoneAreaCode) {
		return xerr.ErrParams.New("请提供有效的手机区号")
	}
	switch req.PhoneAreaCode {
	case "86":
		if len(req.Phone) != 11 {
			return xerr.ErrParams.New("请提供有效的手机号")
		}
	}
	if len(req.VerifyCode) != 4 {
		return xerr.ErrParams.New("请提供有效的验证码")
	}
	// TODO: check req.VerifyCode
	return nil
}
