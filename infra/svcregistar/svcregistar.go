package svcregistar

import "microsvc/deploy"

func Init(must bool) func(*deploy.XConfig, func(must bool, err error)) {
	// 在这里 决定使用 etcd/consul
	return func(cc *deploy.XConfig, onEnd func(must bool, err error)) {
		// todo
		var err error
		onEnd(must, err)
	}
}
