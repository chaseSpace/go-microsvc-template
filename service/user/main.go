package main

import (
	"microsvc/deploy"
	"microsvc/infra"
	"microsvc/infra/cache"
	"microsvc/infra/orm"
	"microsvc/infra/svccli"
	"microsvc/infra/svcregistar"
)

func main() {
	_ = deploy.XConf // init config

	infra.MustSetup(
		cache.InitRedis(true),
		orm.InitGorm(true),
		svcregistar.Init(true),
		svccli.Init(true),
	)
}
