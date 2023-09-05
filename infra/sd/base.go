package sd

import (
	"fmt"
	"go.uber.org/zap"
	"microsvc/deploy"
	"microsvc/infra/sd/abstract"
	"microsvc/infra/sd/consul"
	"microsvc/infra/sd/mdns"
	"microsvc/pkg/xerr"
	"microsvc/pkg/xlog"
	"microsvc/util/ip"
)

var registeredSvc []string

const logPrefix = "sd: "

var rootSD abstract.ServiceDiscovery

func Init(must bool) func(*deploy.XConfig, func(must bool, err error)) {
	return func(cc *deploy.XConfig, onEnd func(must bool, err error)) {
		var err error
		if cc.IsDevEnv() {
			rootSD = mdns.New()
		} else {
			// 在这里 决定使用 etcd/consul
			rootSD, err = NewConsulSD()
			if err != nil {
				xlog.Error(logPrefix+"New failed", zap.Error(err))
			}
		}
		onEnd(must, err)
	}
}

func NewConsulSD() (abstract.ServiceDiscovery, error) {
	cli, err := consul.New()
	if err == nil {
		rootSD = cli
	} else {
		return nil, xerr.ErrInternal.NewMsg(err.Error())
	}
	return cli, nil
}

func Register(reg ...deploy.RegisterSvc) {
	selfIp := "127.0.0.1"
	if !deploy.XConf.IsDevEnv() {
		localIps, err := ip.GetLocalPrivateIPs(true, "")
		if err != nil || len(localIps) == 0 {
			xlog.Panic(logPrefix+"GetLocalPrivateIPs failed", zap.Error(err))
		}
		selfIp = localIps[0].String()
	}

	for _, r := range reg {
		name, addr, port := r.RegGRPCBase()
		if name == "" {
			panic(fmt.Sprintf(logPrefix + "svc'name cannot be empty"))
		}
		if addr == "" {
			addr = selfIp
		}
		err := rootSD.Register(name, addr, port, r.RegGRPCMeta())
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
