package user

import (
	"context"
	"microsvc/enums"
	"microsvc/infra/svccli/rpcext"
	"microsvc/pkg/xerr"
	"microsvc/protocol/svc/user"
	deploy2 "microsvc/service/user/deploy"
	"microsvc/test/tbase"
	"testing"
)

func init() {
	tbase.TearUp(enums.SvcUser, deploy2.UserConf)
}

func TestGetUser(t *testing.T) {
	rsp, err := rpcext.User().GetUser(context.TODO(), &user.GetUserReq{
		Base: nil,
		Uids: nil,
	})
	if !xerr.ErrParams.Equal(err) {
		t.Fatalf("case 1: err is not ErrParams: %v", err)
	}

	rsp, err = rpcext.User().GetUser(context.TODO(), &user.GetUserReq{Uids: []int64{1}})
	if err != nil {
		t.Fatalf("case 2: err %v", err)
	}
	if len(rsp.Umap) != 1 {
		t.Fatalf("case 2: rsp umap %+v", rsp.Umap)
	}
}
