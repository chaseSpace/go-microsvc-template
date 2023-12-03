//go:build k8s

package svccli

import (
	"microsvc/deploy"
	"microsvc/enums/svc"
	"microsvc/infra/sd"
	"microsvc/infra/sd/abstract"
)

func SetDefaultSD(sd abstract.ServiceDiscovery) {}

func Init(must bool) func(*deploy.XConfig, func(must bool, err error)) {
	return func(cc *deploy.XConfig, onEnd func(must bool, err error)) {
		// 不需要实现
	}
}

type RpcClient struct{}

func (c *RpcClient) Getter() any                     { return nil }
func NewCli(svc svc.Svc, gc sd.GenClient) *RpcClient { return nil }
func Stop()                                          {}
