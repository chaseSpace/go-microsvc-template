//go:build k8s

package svccli

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"microsvc/enums/svc"
	"microsvc/infra/sd"
	"microsvc/infra/xgrpc"
	"microsvc/pkg/xlog"
)

func GetConn(svc svc.Svc) (conn *grpc.ClientConn) {
	target := sd.GetSvcTargetInK8s(svc)
	conn, err := xgrpc.NewGRPCClient(target, svc.Name())
	if err != nil {
		xlog.Error("getGRPCClient", zap.Error(err))
		return xgrpc.NewInvalidGRPCConn(svc.Name())
	}
	return conn
}
