package pkg

import "microsvc/deploy"

func Setup(initializers ...deploy.Initializer) {
	for _, initFn := range initializers {
		initFn(deploy.XConf)
	}
}
