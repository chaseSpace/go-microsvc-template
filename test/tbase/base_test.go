package tbase

import (
	"context"
	"github.com/hashicorp/go-uuid"
	"github.com/segmentio/ksuid"
	"microsvc/enums/svc"
	"microsvc/infra/svccli/rpcext"
	"microsvc/pkg/xerr"
	"microsvc/protocol/svc/admin"
	"microsvc/protocol/svc/user"
	deploy2 "microsvc/service/user/deploy"
	"testing"
)

func TestNoRPCClient(t *testing.T) {
	TearUpWithEmptySD(svc.User, deploy2.UserConf)
	defer TearDown()

	_, err := rpcext.Admin().GetUser(context.TODO(), &admin.GetUserReq{
		Base: nil,
		Uids: nil,
	})
	if !xerr.ErrServiceUnavailable.Is(err) {
		t.Errorf("case 1: err is not ErrServiceUnavailable: %v", err)
	}
}

// Run user svc first.
func TestHaveRPCClient(t *testing.T) {
	TearUp(svc.User, deploy2.UserConf)
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

func TestUUID(t *testing.T) {
	imap := make(map[string]interface{})
	for i := 0; i < 10000000; i++ {
		s, err := uuid.GenerateUUID()
		if err != nil {
			t.Fatalf(err.Error())
		}
		if imap[s] != nil {
			t.Fatal("exists", i)
		}
		imap[s] = 1
		println(i, s)
	}
}

func TestKsuid(t *testing.T) {
	imap := make(map[string]interface{})
	// https://github.com/segmentio/ksuid
	// 生成  一种可按生成时间排序、固定20 bytes的 唯一id；无碰撞、无协调、无依赖
	// - 按时间戳按字典顺序排序
	// - base62 编码的文本表示，url友好，复制友好

	s := ksuid.New()
	for i := 0; i < 5000000; i++ {
		s = s.Next()
		if imap[s.String()] != nil {
			t.Fatal("exists", i)
		}
		imap[s.String()] = 1
		println(i, s.String())
	}
}
