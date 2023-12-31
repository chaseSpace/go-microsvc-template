package cache

import (
	"context"
	"fmt"
	"github.com/k0kubun/pp"
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
			fmt.Println("#### infra.redis init success")
			err = setupSvcDB()
			if err != nil {
				panic(err)
			}
		} else {
			pp.Printf("#### infra.redis init failed: %v\n", err)
		}

		onEnd(must, err)
	}
}

type RedisObj struct {
	name deploy.DBname
	*redis.Client
	// 这里可以添加一些其他自定义成员
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
	if len(instMap) > 0 {
		xlog.Debug("cache-redis: resource released...")
	}
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

func IgnoreNil(err error) (bool, error) {
	if err == redis.Nil {
		return true, nil
	}
	return false, err
}
