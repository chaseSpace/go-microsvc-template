//go:build k8s

package sd

import "microsvc/deploy"

const logPrefix = "sd: "

func Init(must bool) func(*deploy.XConfig, func(must bool, err error)) {
	return func(cc *deploy.XConfig, onEnd func(must bool, err error)) {
	}
}

func Stop() {}

func MustRegister(reg ...deploy.RegisterSvc) {}
