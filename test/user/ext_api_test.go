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
测试RPC接口，需要先在本地启动待测试的微服务（可不用启动网关）
*/

func TestHealthCheck(t *testing.T) {
	tbase.GRPCHealthCheck(t, svc.User, deploy2.UserConf)
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
				Phone:         "18855556666", // 若提示已被注册，先删除db账号记录
				VerifyCode:    "xsd1",
			},
			wantErr: nil,
		},
		{
			title: "6.手机号已注册",
			req: &user.SignUpReq{
				Nickname:      "user1", // 昵称限长10字符
				Sex:           enums.SexMale.Int32(),
				Birthday:      "2023-01-01",
				PhoneAreaCode: "86",
				Phone:         "18855556666",
				VerifyCode:    "xsd1",
			},
			wantErr: xerr.ErrParams.New("该手机号已被注册"),
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

	phone := 18855557777
	for i := 0; i < 100; i++ {
		req := &user.SignUpReq{
			Nickname:      fmt.Sprintf("user%d", i), // 昵称限长10字符
			Sex:           enums.SexMale.Int32(),
			Birthday:      "2023-01-01",
			PhoneAreaCode: "86",
			Phone:         fmt.Sprintf("%d", phone+i),
			VerifyCode:    "xzdw",
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

	uidRepeatedErrCnt := atomic.Int32{}

	expectedErr := xerr.ErrInternal.New("太多人注册辣，隔几秒再试一下哦")
	// 并发注册

	// 业务代码中设置的号池容量为100，这里设置相同的并发参数100，则完全可以处理，不会出现重复ID
	// - 对于更高的并发，虽然也不会出现重复id，但会增加接口耗时，同时建议提高 号池容量 配置以提高并发性能
	total := 100
	phone := 18855560818
	for i := 0; i < total; i++ {
		req := &user.SignUpReq{
			Nickname:      fmt.Sprintf("user%d", i), // 昵称限长10字符
			Sex:           enums.SexMale.Int32(),
			Birthday:      "2023-01-01",
			PhoneAreaCode: "86",
			Phone:         fmt.Sprintf("%d", phone+i),
			VerifyCode:    "xzdw",
		}
		if i%2 == 0 {
			req.Sex = enums.SexFemale.Int32()
		}
		x.Add(1)
		go func() {
			_, err := rpcext.User().SignUp(tbase.NewTestCallCtx(), req)
			if err != nil {
				assert.Equal(t, expectedErr, err)
				uidRepeatedErrCnt.Add(1)
			}
			x.Done()
		}()
	}

	x.Wait()

	errCnt := uidRepeatedErrCnt.Load()
	if errCnt > 0 {
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

func TestGetUser(t *testing.T) {
	tbase.TearUp(svc.User, deploy2.UserConf)
	defer tbase.TearDown()

	// case-1: need `base` arg
	_, err := rpcext.User().GetUser(tbase.TestCallCtx, &user.GetUserReq{
		Base: nil,
		Uids: nil,
	})
	assert.Equal(t, xerr.ErrParams.AppendMsg("missing arg:`base`"), err)

	// case-2: need `uids` arg
	_, err = rpcext.User().GetUser(tbase.TestCallCtx, &user.GetUserReq{
		Base: tbase.TestBaseExtReq,
		Uids: nil,
	})
	assert.Equal(t, xerr.ErrParams.New("missing arg:`uids`"), err)

	// case-3: need valid `uids` arg (at least exist one)
	_, err = rpcext.User().GetUser(tbase.TestCallCtx, &user.GetUserReq{
		Base: tbase.TestBaseExtReq,
		Uids: []int64{0, 0, 0},
	})
	assert.Equal(t, xerr.ErrParams.New("missing arg:`uids`"), err)

	// case-4: normal（如果 id:100010 不存在，请先运行上面的用例 TestSignUp 来注册id）
	rsp, err := rpcext.User().GetUser(tbase.TestCallCtx, &user.GetUserReq{
		Base: tbase.TestBaseExtReq,
		Uids: []int64{100010, 0, 999}, // 仅 100010 有效
	})
	assert.Nil(t, err)

	// 只返回有效的uid
	expectedMap := map[int64]*user.User{
		100010: &user.User{
			Uid:      100010,
			Nickname: "user1",
			Birthday: "2023-01-01",
			Sex:      1,
		}}
	assert.Equal(t, len(expectedMap), len(rsp.Umap))
	for key, val := range expectedMap {
		actualVal, ok := rsp.Umap[key]
		assert.True(t, ok, "Key不存在于实际的map中: %s", key)
		assert.EqualExportedValues(t, *val, *actualVal)
		//assert.Truef(t, proto.Equal(val, actualVal), "Val不相等: actualVal -> %+v", actualVal)
	}
}
