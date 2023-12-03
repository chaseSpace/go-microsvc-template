//go:build !k8s

package xgrpc

import (
	"microsvc/deploy"
	"microsvc/util"
	"net"
)

func getListener(defaultPort int, portSetter deploy.SvcListenPortSetter) (net.Listener, int, error) {
	lisFetcher := util.NewTcpListenerFetcher(grpcPortMin, grpcPortMax)
	lis, port, err := lisFetcher.Get()
	if err != nil {
		return nil, 0, err
	}
	portSetter.SetGRPC(port)
	return lis, port, err
}
