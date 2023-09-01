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
	"sync"
)

var oncemap sync.Map

func TearUp(svc enums.Svc, svcConf deploy.SvcConfImpl) {
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

func TearUpWithEmptySD(svc enums.Svc, svcConf deploy.SvcConfImpl) {
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
		svccli.SetDefaultSD(abstract.EmptySD{})
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
