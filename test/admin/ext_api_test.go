package user

import (
	"context"
	"microsvc/enums"
	"microsvc/infra/svccli/rpcext"
	"microsvc/pkg/xerr"
	"microsvc/protocol/svc/admin"
	"microsvc/service/admin/deploy"
	"microsvc/test/tbase"
	"testing"
)

func init() {
	tbase.TearUp(enums.SvcAdmin, deploy.AdminConf)
}

func TestGetUser(t *testing.T) {
	rsp, err := rpcext.Admin().GetUser(context.TODO(), &admin.GetUserReq{
		Base: nil,
		Uids: nil,
	})
	if !xerr.ErrParams.Equal(err) {
		t.Fatalf("case 1: err is not ErrParams: %v", err)
	}

	rsp, err = rpcext.Admin().GetUser(context.TODO(), &admin.GetUserReq{
		Base: nil,
		Uids: []int64{1},
	})
	if err != nil {
		t.Fatal("case 2: err should not be nil")
	}
	if len(rsp.Umap) != 1 {
		t.Fatalf("case 2: err rsp.umap")
	}
}
