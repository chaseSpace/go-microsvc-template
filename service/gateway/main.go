package main

import (
	"microsvc/deploy"
	"microsvc/enums"
	"microsvc/infra/xhttp"
	deploy2 "microsvc/service/gateway/deploy"
	"microsvc/service/gateway/handler"
	"microsvc/util/graceful"
)

func main() {
	graceful.SetupSignal()
	defer graceful.OnExit()

	deploy.Init(enums.SvcGateway, deploy2.GatewayConf)

	ctrl := new(handler.GatewayCtrl)
	server := xhttp.New(deploy2.GatewayConf.HttpPort, ctrl.Handler)

	graceful.Schedule(server.Start)
	graceful.Run()
}
