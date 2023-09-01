package main

import (
	"microsvc/deploy"
	"microsvc/enums"
	"microsvc/infra"
	"microsvc/infra/svccli"
	_ "microsvc/infra/xgrpc/proto"
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

	deploy.Init(enums.SvcGateway, deploy2.GatewayConf)

	pkg.Setup(
		xlog.Init,
		xkafka.Init,
	)
	infra.Setup(
		svccli.Init(true),
	)

	ctrl := new(handler.GatewayCtrl)
	server := xhttp.New(deploy2.GatewayConf.HttpPort, ctrl.Handler)

	graceful.Schedule(server.Start)
	graceful.Run()
}
