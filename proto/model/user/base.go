package user

import (
	"microsvc/infra/cache"
	"microsvc/infra/orm"
	"microsvc/proto/model"
)

var (
	Q    = orm.NewMysqlObj(model.MysqlDB)
	QLog = orm.NewMysqlObj(model.MysqlDBLog)
)

var (
	R    = cache.NewRedisObj(model.RedisDB)
	RAdm = cache.NewRedisObj(model.RedisDBAdmin)
)

func init() {
	// 此函数会在main函数执行前向orm注入服务需要使用的DB对象
	orm.Setup(Q, QLog)
	cache.Setup(R)
}
