package tbase

import (
	"context"
	"microsvc/enums/svc"
	"microsvc/infra/svccli/rpcext"
	"microsvc/pkg/xerr"
	"microsvc/protocol/svc/admin"
	"microsvc/protocol/svc/user"
	deploy2 "microsvc/service/user/deploy"
	"testing"
)

func TestNoRPCClient(t *testing.T) {
	TearUpWithEmptySD(svc.SvcUser, deploy2.UserConf)
	defer TearDown()

	_, err := rpcext.Admin().GetUser(context.TODO(), &admin.GetUserReq{
		Base: nil,
		Uids: nil,
	})
	if !xerr.ErrNoRPCClient.Is(err) {
		t.Errorf("case 1: err is not ErrNoRPCClient: %v", err)
	}
}

// Run user svc first.
func TestHaveRPCClient(t *testing.T) {
	TearUp(svc.SvcUser, deploy2.UserConf)
	defer TearDown()

	rsp, err := rpcext.User().GetUser(context.TODO(), &user.GetUserReq{
		Base: nil,
		Uids: []int64{1},
	})
	if err != nil {
		t.Errorf("case 1: err: %v", err)
	} else {
		if len(rsp.Umap) != 1 {
			t.Errorf("case 1: err rsp: %+v", rsp.Umap)
		}
	}
}
