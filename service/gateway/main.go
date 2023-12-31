package main

import (
	"microsvc/deploy"
	"microsvc/enums/svc"
	"microsvc/infra"
	"microsvc/infra/svccli"
	_ "microsvc/infra/xgrpc/protobytes"
	"microsvc/infra/xhttp"
	"microsvc/pkg"
	"microsvc/pkg/xkafka"
	"microsvc/pkg/xlog"
	deploy2 "microsvc/service/gateway/deploy"
	"microsvc/service/gateway/handler"
	"microsvc/util/graceful"
)

func main() {
	graceful.SetupSignal()
	defer graceful.OnExit()

	deploy.Init(svc.Gateway, deploy2.GatewayConf)

	pkg.Setup(
		xlog.Init,
		xkafka.Init,
	)
	infra.Setup(
		svccli.Init(true),
	)

	ctrl := new(handler.GatewayCtrl)
	server := xhttp.New(deploy2.GatewayConf.HttpPort, ctrl.Handler)

	graceful.Register(server.Start)
	graceful.Run()
}
