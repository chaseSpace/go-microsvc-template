package svcdiscovery

import (
	"github.com/pkg/errors"
	"microsvc/deploy"
	"microsvc/infra/svcdiscovery/consul"
	"microsvc/infra/svcdiscovery/define"
)

func Init(must bool) func(*deploy.XConfig, func(must bool, err error)) {
	return func(cc *deploy.XConfig, onEnd func(must bool, err error)) {
		// 在这里 决定使用 etcd/consul
		cli, err := consul.NewConsulSD()
		if err == nil {
			Sd = cli
		} else {
			err = errors.Wrap(err, "NewConsulSD")
		}
		onEnd(must, err)
	}
}

var Sd define.ServiceDiscovery
