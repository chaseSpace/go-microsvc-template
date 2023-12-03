//go:build k8s

package rpcext

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"microsvc/enums/svc"
	"microsvc/infra/sd"
	"microsvc/infra/xgrpc"
	"microsvc/pkg/xlog"
	"microsvc/protocol/svc/review"
	"microsvc/protocol/svc/user"
)

// If you use this file, service client is directly use DNS name as service target address (e.g. in K8s environment).

func getGRPCClient(svc svc.Svc) *grpc.ClientConn {
	// e.g. `go-svc-user` is `user` service's dns name.
	target := sd.GetSvcTargetInK8s(svc)
	conn, err := xgrpc.NewGRPCClient(target, svc.Name())
	if err != nil {
		xlog.Error("getGRPCClient", zap.Error(err))
		return xgrpc.NewInvalidGRPCConn(svc.Name())
	}
	return conn
}

func User() user.UserExtClient {
	return user.NewUserExtClient(getGRPCClient(svc.User))
}

func Review() review.ReviewExtClient {
	return review.NewReviewExtClient(getGRPCClient(svc.Review))
}
