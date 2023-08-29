package user

import (
	"context"
	"microsvc/enums"
	"microsvc/infra/svccli/rpc"
	"microsvc/protocol/svc/user"
	deploy2 "microsvc/service/user/deploy"
	"microsvc/test/tbase"
	"testing"
)

func TestGetUserInt(t *testing.T) {
	tbase.TearUp(enums.SvcUser, deploy2.UserConf)
	defer tbase.TearDown()

	rsp, err := rpc.User().GetUser(context.TODO(), &user.GetUserIntReq{Uids: []int64{1}})
	if err != nil {
		t.Fatal("err", err)
	}
	t.Logf("rsp: %+v", rsp)

	rsp, err = rpc.User().GetUser(context.TODO(), &user.GetUserIntReq{})
	if err == nil {
		t.Fatal("case 2: err should not be nil")
	}
	t.Logf("error is right: %v", err)
}
