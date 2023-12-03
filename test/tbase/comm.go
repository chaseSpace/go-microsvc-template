package tbase

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"microsvc/consts"
	"microsvc/deploy"
	"microsvc/enums/svc"
	"microsvc/infra"
	"microsvc/infra/sd/abstract"
	"microsvc/infra/svccli"
	"microsvc/infra/xgrpc"
	"microsvc/pkg"
	"microsvc/pkg/xlog"
	svc2 "microsvc/protocol/svc"
	"microsvc/util"
	"microsvc/util/graceful"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

var oncemap sync.Map

// TearUp 这个方法只完成grpc服务的client初始化，所以你需要提前在本地启动待测试的服务
func TearUp(svc svc.Svc, svcConf deploy.SvcConfImpl) {
	var o = new(sync.Once)
	v, ok := oncemap.Load(svc.Name())
	if !ok {
		oncemap.Store(svc.Name(), o)
	} else {
		o = v.(*sync.Once)
	}
	o.Do(func() {
		//println(111, svc.Name())
		_ = os.Setenv(consts.EnvVarLogLevel, "debug")
		graceful.SetupSignal()

		if !isProjectRootDir() {
			wd, _ := os.Getwd()
			parentDir := filepath.Dir(filepath.Dir(wd))
			_ = os.Chdir(parentDir)
		}
		deploy.Init(svc, svcConf)
		pkg.Setup(
			xlog.Init,
		)
		infra.Setup(
			//sd.Init(true),
			svccli.Init(true),
		)
	})
}

func TearDown() {
	graceful.OnExit()
}

var oncemapEmptySD sync.Map

func TearUpWithEmptySD(svc svc.Svc, svcConf deploy.SvcConfImpl) {
	var o = new(sync.Once)
	v, ok := oncemapEmptySD.Load(svc.Name())
	if !ok {
		oncemapEmptySD.Store(svc.Name(), o)
	} else {
		o = v.(*sync.Once)
	}

	o.Do(func() {
		_ = os.Setenv(consts.EnvVarLogLevel, "debug")
		graceful.SetupSignal()

		if !isProjectRootDir() {
			wd, _ := os.Getwd()
			parentDir := filepath.Dir(filepath.Dir(wd))
			_ = os.Chdir(parentDir)
		}
		deploy.Init(svc, svcConf)
		pkg.Setup(
			xlog.Init,
		)
		svccli.SetDefaultSD(abstract.Empty{})
		infra.Setup(
			//sd.Init(true),
			svccli.Init(true),
		)
	})
}

func isProjectRootDir() bool {
	_, err := os.Stat("go.mod")
	return err == nil
}

// TestCallCtx 它的traceId在多次使用时是同一个，若需要不同的traceId请使用 NewTestCallCtx
var TestCallCtx = NewTestCallCtx()

func NewTestCallCtx() context.Context {
	md := metadata.Pairs(
		xgrpc.MdKeyTestCall, xgrpc.MdKeyFlagExist,
		xgrpc.MdKeyTraceId, util.NewKsuid(),
	)
	return metadata.NewOutgoingContext(context.TODO(), md)
}

var TestBaseExtReq = &svc2.BaseExtReq{
	App:        "test_app",
	AppVersion: "1.0.0",
	Extension:  nil,
}

var TestBaseAdminReq = &svc2.AdminBaseReq{
	Uid:       1,
	Nick:      "Lucy",
	Extension: nil,
}

func GRPCHealthCheck(t *testing.T, svc2 svc.Svc, svcConf deploy.SvcConfImpl) {
	TearUp(svc2, svcConf)
	defer TearDown()

	healthCli := svccli.NewCli(svc2, func(conn *grpc.ClientConn) interface{} { return grpc_health_v1.NewHealthClient(conn) })
	cli := healthCli.Getter().(grpc_health_v1.HealthClient)

	response, err := cli.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{
		Service: svc2.Name(),
	})
	if err != nil {
		panic(err)
	}

	assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, response.Status)
}
