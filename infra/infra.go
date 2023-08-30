package infra

import (
	"go.uber.org/zap"
	"microsvc/deploy"
	"microsvc/infra/cache"
	"microsvc/infra/orm"
	"microsvc/infra/sd"
	"microsvc/infra/svccli"
	"microsvc/pkg/xlog"
	"microsvc/util/graceful"
)

type initFunc func(cc *deploy.XConfig, onEnd func(must bool, err error))

func Setup(initFn ...initFunc) {
	for _, fn := range initFn {
		fn(deploy.XConf, func(must bool, err error) {
			if must && err != nil {
				panic(err)
			}
			if err != nil {
				xlog.Error("infra.Setup err", zap.Error(err))
			}
		})
	}

	graceful.AddStopFunc(func() {
		svccli.Stop()
		sd.Stop()
		cache.Stop()
		orm.Stop()
	})
}
