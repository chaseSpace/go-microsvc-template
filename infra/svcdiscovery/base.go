package svcdiscovery

import (
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"microsvc/deploy"
	"microsvc/infra/svcdiscovery/consul"
	"microsvc/infra/svcdiscovery/sd"
	"microsvc/pkg/xlog"
	"microsvc/util/ip"
)

var registeredSvc []string

const logPrefix = "svcdiscovery: "

var rootSD sd.ServiceDiscovery

func Init(must bool, reg ...deploy.RegisterSvc) func(*deploy.XConfig, func(must bool, err error)) {
	return func(cc *deploy.XConfig, onEnd func(must bool, err error)) {
		// 在这里 决定使用 etcd/consul
		cli, err := consul.NewConsulSD()
		if err == nil {
			rootSD = cli
			localIps, err := ip.GetLocalPrivateIPs(true, "")
			if err != nil || len(localIps) == 0 {
				xlog.Error(logPrefix+"GetLocalPrivateIPs failed, stop register", zap.Error(err), zap.Int("ip len", len(localIps)))
			} else {
				if len(reg) == 0 && deploy.XConf.GetSvcConf() != nil {
					reg = append(reg, deploy.XConf.GetSvcConf())
				}
				for _, r := range reg {
					name, addr, port := r.RegBase()
					if name == "" {
						panic(fmt.Sprintf(logPrefix + "svc'name cannot be empty"))
					}
					if addr == "" {
						addr = localIps[0].String()
					}
					err = rootSD.Register(name, addr, port, r.RegMeta())
					if err != nil {
						xlog.Error(logPrefix+"Register Svc failed, stop register", zap.String("Svc", name), zap.Error(err))
						break
					}
					xlog.Error(logPrefix+"Register Svc success", zap.String("reg_svc", name), zap.String("addr", fmt.Sprintf("%s:%d", addr, port)))
					registeredSvc = append(registeredSvc, name)
				}
			}
		} else {
			xlog.Error(logPrefix+"NewSD failed", zap.Error(err))
			err = errors.Wrap(err, "NewSD")
		}
		onEnd(must, err)
	}
}

func GetSD() sd.ServiceDiscovery {
	return rootSD
}

func Stop() {
	for _, s := range registeredSvc {
		err := rootSD.Deregister(s)
		println(111, s)
		if err != nil {
			xlog.Error(logPrefix+"Deregister fail", zap.Error(err))
		}
	}
}
