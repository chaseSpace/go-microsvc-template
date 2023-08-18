package pkg

import "microsvc/deploy"

func Init(initializers ...deploy.Initializer) {
	for _, initFn := range initializers {
		initFn(deploy.XConf)
	}
}
