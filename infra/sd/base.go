package sd

import (
	"fmt"
	"go.uber.org/zap"
	"microsvc/deploy"
	"microsvc/infra/sd/abstract"
	"microsvc/infra/sd/consul"
	"microsvc/pkg/xerr"
	"microsvc/pkg/xlog"
	"microsvc/util/ip"
)

var registeredSvc []string

const logPrefix = "sd: "

var rootSD abstract.ServiceDiscovery

func Init(must bool) func(*deploy.XConfig, func(must bool, err error)) {
	return func(cc *deploy.XConfig, onEnd func(must bool, err error)) {
		// 在这里 决定使用 etcd/consul
		cli, err := NewSD()
		if err != nil {
			xlog.Error(logPrefix+"NewSD failed", zap.Error(err))
		} else {
			rootSD = cli
		}
		onEnd(must, err)
	}
}

func NewSD() (abstract.ServiceDiscovery, error) {
	cli, err := consul.NewConsulSD()
	if err == nil {
		rootSD = cli
	} else {
		return nil, xerr.ErrInternal.NewMsg(err.Error())
	}
	return cli, nil
}

func Register(reg ...deploy.RegisterSvc) {
	localIps, err := ip.GetLocalPrivateIPs(true, "")
	if err != nil || len(localIps) == 0 {
		xlog.Panic(logPrefix+"GetLocalPrivateIPs failed", zap.Error(err))
	}
	for _, r := range reg {
		name, addr, port := r.RegGRPCBase()
		if name == "" {
			panic(fmt.Sprintf(logPrefix + "svc'name cannot be empty"))
		}
		if addr == "" {
			addr = localIps[0].String()
		}
		err = rootSD.Register(name, addr, port, r.RegGRPCMeta())
		if err != nil {
			xlog.Error(logPrefix+"register svc failed", zap.String("Svc", name), zap.Error(err))
			break
		}
		xlog.Info(logPrefix+"register svc success", zap.String("reg_svc", name), zap.String("addr", fmt.Sprintf("%s:%d", addr, port)))
		registeredSvc = append(registeredSvc, name)
	}
}

func Stop() {
	for _, s := range registeredSvc {
		err := rootSD.Deregister(s)
		if err != nil {
			xlog.Error(logPrefix+"deregister fail", zap.Error(err))
		} else {
			xlog.Debug(logPrefix+"deregister success", zap.String("svc", s))
		}
	}
}
