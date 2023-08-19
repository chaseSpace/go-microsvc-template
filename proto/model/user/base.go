package user

import (
	"microsvc/infra/cache"
	"microsvc/infra/orm"
)

const (
	mysqlDBname  = "microsvc"
	mysqlDBname2 = "microsvc_log"
)
const redisDBname = "microsvc"

var (
	Q    = orm.NewMysqlObj(mysqlDBname)
	QLog = orm.NewMysqlObj(mysqlDBname2)
)

var (
	R = cache.NewRedisObj(redisDBname)
)

func init() {
	// 此函数会在main函数执行前向orm注入服务需要使用的DB对象
	orm.RegSvcDB(Q, QLog)
	cache.RegSvcDB(R)
}
