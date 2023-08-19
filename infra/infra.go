package infra

import (
	"microsvc/deploy"
	"microsvc/infra/cache"
	"microsvc/infra/orm"
)

type initFunc func(cc *deploy.XConfig, onEnd func(must bool, err error))

func MustSetup(initFn ...initFunc) {
	for _, fn := range initFn {
		fn(deploy.XConf, func(must bool, err error) {
			if must && err != nil {
				panic(err)
			}
			// TODO LOG
		})
	}
}

func Stop() {
	orm.Stop()
	cache.Stop()
}
