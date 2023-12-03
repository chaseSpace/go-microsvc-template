//go:build k8s

package xgrpc

import (
	"fmt"
	"microsvc/deploy"
	"net"
)

func getListener(defaultPort int, portSetter deploy.SvcListenPortSetter) (net.Listener, int, error) {
	//portSetter.SetGRPC(defaultPort) 使用DNS时不需要
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", defaultPort))
	return lis, defaultPort, err
}
