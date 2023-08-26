package infra

import (
	"go.uber.org/zap"
	"microsvc/deploy"
	"microsvc/infra/cache"
	"microsvc/infra/orm"
	"microsvc/infra/svccli"
	"microsvc/infra/svcdiscovery"
	"microsvc/pkg/xlog"
)

type initFunc func(cc *deploy.XConfig, onEnd func(must bool, err error))

func MustSetup(initFn ...initFunc) {
	for _, fn := range initFn {
		fn(deploy.XConf, func(must bool, err error) {
			if must && err != nil {
				panic(err)
			}
			if err != nil {
				xlog.Error("infra.MustSetup err", zap.Error(err))
			}
		})
	}
}

func Stop() {
	orm.Stop()
	cache.Stop()
	svccli.Stop()
	svcdiscovery.Stop()
}
