package tbase

import (
	"microsvc/consts"
	"microsvc/deploy"
	"microsvc/enums"
	"microsvc/infra"
	"microsvc/infra/sd"
	"microsvc/infra/sd/abstract"
	"microsvc/infra/svccli"
	"microsvc/pkg"
	"microsvc/pkg/xlog"
	"microsvc/util/graceful"
	"os"
	"path/filepath"
)

func TearUp(svc enums.Svc, svcConf deploy.SvcConfImpl) {
	_ = os.Setenv(consts.EnvVarLogLevel, "debug")
	graceful.SetupSignal()

	wd, _ := os.Getwd()
	parentDir := filepath.Dir(filepath.Dir(wd))
	_ = os.Chdir(parentDir)
	deploy.Init(svc, svcConf)
	pkg.Init(
		xlog.Init,
	)
	infra.MustSetup(
		sd.Init(true),
		svccli.Init(true),
	)
}

func TearUpWithEmptySD(svc enums.Svc, svcConf deploy.SvcConfImpl) {
	_ = os.Setenv(consts.EnvVarLogLevel, "debug")
	graceful.SetupSignal()

	wd, _ := os.Getwd()
	parentDir := filepath.Dir(filepath.Dir(wd))
	_ = os.Chdir(parentDir)
	deploy.Init(svc, svcConf)
	pkg.Init(
		xlog.Init,
	)
	svccli.SetDefaultSD(abstract.EmptySD{})
	infra.MustSetup(
		//sd.Init(true),
		svccli.Init(true),
	)
}

func TearDown() {
	graceful.OnExit()
}
