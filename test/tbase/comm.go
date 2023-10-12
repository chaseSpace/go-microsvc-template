package tbase

import (
	"context"
	"google.golang.org/grpc/metadata"
	"microsvc/consts"
	"microsvc/deploy"
	"microsvc/enums/svc"
	"microsvc/infra"
	"microsvc/infra/sd"
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
)

var oncemap sync.Map

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
			sd.Init(true),
			svccli.Init(true),
		)
	})
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
			//sd.Setup(true),
			svccli.Init(true),
		)
	})
}

func TearDown() {
	graceful.OnExit()
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
	ThisIsExtApi: true,
	App:          "test_app",
	AppVersion:   "1.0.0",
	Extension:    nil,
}
