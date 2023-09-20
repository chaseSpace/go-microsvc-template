package user

import (
	"context"
	"microsvc/enums/svc"
	"microsvc/infra/svccli/rpcext"
	"microsvc/pkg/xerr"
	"microsvc/protocol/svc/user"
	deploy2 "microsvc/service/user/deploy"
	"microsvc/test/tbase"
	"testing"
)

func TestGetUser(t *testing.T) {
	tbase.TearUp(svc.User, deploy2.UserConf)
	defer tbase.TearDown()

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
