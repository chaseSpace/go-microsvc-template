package user

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"microsvc/enums"
	"microsvc/enums/svc"
	"microsvc/infra/svccli/rpcext"
	"microsvc/pkg/xerr"
	"microsvc/protocol/svc/user"
	deploy2 "microsvc/service/user/deploy"
	"microsvc/test/tbase"
	"sync"
	"sync/atomic"
	"testing"
)

/*
测试RPC接口，需要先在本地启动待测试的微服务（不用启动网关）
*/

func TestGetUser(t *testing.T) {
	tbase.TearUp(svc.User, deploy2.UserConf)
	defer tbase.TearDown()

	rsp, err := rpcext.User().GetUser(tbase.TestCallCtx, &user.GetUserReq{
		Base: nil,
		Uids: nil,
	})
	if !xerr.ErrParams.Equal(err) {
		t.Fatalf("case 1: err is not ErrParams: %v", err)
	}

	rsp, err = rpcext.User().GetUser(tbase.TestCallCtx, &user.GetUserReq{Uids: []int64{1}})
	if err != nil {
		t.Fatalf("case 2: err %v", err)
	}
	if len(rsp.Umap) != 2 {
		t.Fatalf("case 2: rsp umap %+v", rsp.Umap)
	}
}

func TestSignUp(t *testing.T) {
	tbase.TearUp(svc.User, deploy2.UserConf)
	defer tbase.TearDown()

	type item struct {
		title   string
		req     *user.SignUpReq
		wantErr error
	}
	tt := []item{
		{
			title: "1.无效昵称或超出长度",
			req: &user.SignUpReq{
				Nickname:   "user1234567", // 昵称限长10字符
				Sex:        enums.SexMale.Int32(),
				Birthday:   "2024-01-02",
				VerifyCode: "",
			},
			wantErr: xerr.ErrParams.New("无效昵称或超出长度"),
		},
		{
			title: "2.无效的生日信息",
			req: &user.SignUpReq{
				Nickname:   "user1", // 昵称限长10字符
				Sex:        enums.SexMale.Int32(),
				Birthday:   "",
				VerifyCode: "",
			},
			wantErr: xerr.ErrParams.New("无效的生日信息"),
		},
		{
			title: "3.请设置有效的性别",
			req: &user.SignUpReq{
				Nickname:   "user1", // 昵称限长10字符
				Sex:        enums.SexUnknown.Int32(),
				Birthday:   "2023-01-01",
				VerifyCode: "",
			},
			wantErr: xerr.ErrParams.New("请设置有效的性别"),
		},
		{
			title: "4.请提供有效的手机区号",
			req: &user.SignUpReq{
				Nickname:      "user1", // 昵称限长10字符
				Sex:           enums.SexMale.Int32(),
				Birthday:      "2023-01-01",
				PhoneAreaCode: "",
				Phone:         "",
				VerifyCode:    "xsd1",
			},
			wantErr: xerr.ErrParams.New("请提供有效的手机区号"),
		},
		{
			title: "5.请提供有效的手机号",
			req: &user.SignUpReq{
				Nickname:      "user1", // 昵称限长10字符
				Sex:           enums.SexMale.Int32(),
				Birthday:      "2023-01-01",
				PhoneAreaCode: "86",
				Phone:         "",
				VerifyCode:    "xsd1",
			},
			wantErr: xerr.ErrParams.New("请提供有效的手机号"),
		},
		{
			title: "OK",
			req: &user.SignUpReq{
				Nickname:      "user1", // 昵称限长10字符
				Sex:           enums.SexMale.Int32(),
				Birthday:      "2023-01-01",
				PhoneAreaCode: "86",
				Phone:         "18855556666",
				VerifyCode:    "xsd1",
			},
			wantErr: nil,
		},
	}

	for _, v := range tt {
		_, err := rpcext.User().SignUp(tbase.TestCallCtx, v.req)
		assert.Equal(t, v.wantErr, err)
	}
}

func TestBatchSignUp(t *testing.T) {
	tbase.TearUp(svc.User, deploy2.UserConf)
	defer tbase.TearDown()

	for i := 0; i < 100; i++ {
		req := &user.SignUpReq{
			Nickname:   fmt.Sprintf("user%d", i), // 昵称限长10字符
			Sex:        enums.SexMale.Int32(),
			Birthday:   "2023-01-01",
			VerifyCode: "",
		}
		if i%2 == 0 {
			req.Sex = enums.SexFemale.Int32()
		}
		_, err := rpcext.User().SignUp(tbase.TestCallCtx, req)
		assert.Equal(t, nil, err)
	}
}

func TestConcurrencySignUp(t *testing.T) {
	tbase.TearUp(svc.User, deploy2.UserConf)
	defer tbase.TearDown()

	x := sync.WaitGroup{}

	genUIDRepeatedErrCnt := atomic.Int32{}

	expectedErr := xerr.ErrInternal.New("太多人注册辣，隔几秒再试一下哦")
	// 并发注册
	total := 100
	for i := 0; i < total; i++ {
		req := &user.SignUpReq{
			Nickname:   fmt.Sprintf("user%d", i), // 昵称限长10字符
			Sex:        enums.SexMale.Int32(),
			Birthday:   "2023-01-01",
			VerifyCode: "",
		}
		if i%2 == 0 {
			req.Sex = enums.SexFemale.Int32()
		}
		x.Add(1)
		go func() {
			_, err := rpcext.User().SignUp(tbase.NewTestCallCtx(), req)
			if err != nil {
				assert.Equal(t, expectedErr, err)
				genUIDRepeatedErrCnt.Add(1)
			}
			x.Done()
		}()
	}

	x.Wait()

	errCnt := genUIDRepeatedErrCnt.Load()
	if errCnt > 40 {
		t.Errorf("并发次数：%d, 失败次数:%d 超出预期\n", total, errCnt)
	} else {
		t.Logf("并发次数：%d, 失败次数:%d 符合预期\n", total, errCnt)
	}
}

func TestSignIn(t *testing.T) {
	tbase.TearUp(svc.User, deploy2.UserConf)
	defer tbase.TearDown()

	type item struct {
		title   string
		req     *user.SignInReq
		wantErr error
	}
	tt := []item{
		{
			title: "无效手机区号",
			req: &user.SignInReq{
				Base:          tbase.TestBaseExtReq,
				PhoneAreaCode: "",
				Phone:         "",
				VerifyCode:    "",
			},
			wantErr: xerr.ErrParams.New("请提供有效的手机区号"),
		},
		{
			title: "无效手机号",
			req: &user.SignInReq{
				Base:          tbase.TestBaseExtReq,
				PhoneAreaCode: "86",
				Phone:         "",
				VerifyCode:    "",
			},
			wantErr: xerr.ErrParams.New("请提供有效的手机号"),
		},
		{
			title: "无效验证码",
			req: &user.SignInReq{
				Base:          tbase.TestBaseExtReq,
				PhoneAreaCode: "86",
				Phone:         "18855556666",
				VerifyCode:    "",
			},
			wantErr: xerr.ErrParams.New("请提供有效的验证码"),
		},
		{
			title: "手机号未注册",
			req: &user.SignInReq{
				Base:          tbase.TestBaseExtReq,
				PhoneAreaCode: "86",
				Phone:         "0123456789x",
				VerifyCode:    "sd81",
			},
			wantErr: xerr.ErrParams.New("手机号未注册"),
		},
		{
			title: "OK（此case通过 需要先调用上面的 TestSignUp ）",
			req: &user.SignInReq{
				Base:          tbase.TestBaseExtReq,
				PhoneAreaCode: "86",
				Phone:         "18855556666",
				VerifyCode:    "sd81",
			},
			wantErr: nil,
		},
	}

	for _, v := range tt {
		r, err := rpcext.User().SignIn(tbase.TestCallCtx, v.req)
		assert.Equal(t, v.wantErr, err)

		if v.wantErr == nil {
			assert.NotEmpty(t, r.Token)
			assert.NotEmpty(t, r.Info)
		}
	}
}
