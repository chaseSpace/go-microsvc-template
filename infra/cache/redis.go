package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"microsvc/deploy"
	"microsvc/pkg/xlog"
	"microsvc/util"
	"time"
)

var instMap = make(map[deploy.DBname]*redis.Client)

func InitRedis(must bool) func(*deploy.XConfig, func(must bool, err error)) {
	return func(cc *deploy.XConfig, onEnd func(must bool, err error)) {

		var err error
		for _, v := range cc.Redis {
			rdb := redis.NewClient(&redis.Options{
				Addr:       v.Addr,
				Password:   v.Password,
				DB:         v.DB,
				MaxRetries: 2,
			})
			util.RunTaskWithCtxTimeout(time.Second, func(ctx context.Context) {
				err = rdb.Ping(ctx).Err()
			})
			if err != nil {
				break
			}
			instMap[v.DBname] = rdb
		}

		if err == nil {
			err = setupSvcDB()
		}

		onEnd(must, err)
	}
}

type RedisObj struct {
	name deploy.DBname
	*redis.Client
	// 你可能希望在对象中包含一些其他自定义成员，在这里添加
}

func (m *RedisObj) IsInvalid() bool {
	return m.Client == nil
}

func (m *RedisObj) Stop() {
	err := m.Client.Close()
	if err != nil {
		xlog.Error("orm.Stop() failed", zap.Error(err))
	}
}

func (m *RedisObj) String() string {
	return fmt.Sprintf("RedisObj{name:%s, instExists:%v}", m.name, m.Client != nil)
}

var servicesDB []*RedisObj

func setupSvcDB() error {
	for _, obj := range servicesDB {
		obj.Client = instMap[obj.name]
		if obj.IsInvalid() {
			return fmt.Errorf("cache.RedisObj is invalid, %s", obj)
		}
	}
	return nil
}

func Stop() {
	for _, db := range instMap {
		_ = db.Close()
	}
	xlog.Debug("cache-redis: resource released...")
}

func NewRedisObj(dbname deploy.DBname) *RedisObj {
	o := &RedisObj{name: dbname}
	return o
}

func Setup(obj ...*RedisObj) {
	for _, o := range obj {
		if o.name == "" {
			panic(fmt.Sprintf("cache.AddSvcDB: need name"))
		}
	}
	servicesDB = obj
}
