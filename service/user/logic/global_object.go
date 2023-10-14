package logic

import (
	"fmt"
	"github.com/dlclark/regexp2"
	"microsvc/infra/cache"
	"microsvc/proto/model"
	cache2 "microsvc/service/user/cache"
	"microsvc/service/user/dao"
	"microsvc/util"
	"microsvc/util/xlock"
	"microsvc/xvendor/genuserid2"
)

// 一些需要全局使用的资源在这里初始化
var (
	uidGenerator genuserid2.UIDGeneratorApi
)

func MustInit() {
	g := globalObjectCtrl{}
	util.AssertNilErr(g.InitUidGenerator())
}

type globalObjectCtrl struct {
}

func (globalObjectCtrl) InitUidGenerator() error {
	skipPattern := []string{
		`(\d)\1(\d)\2$`, // aabb结尾模式
		`(\d)\1{2}$`,    // aaa结尾模式，包含3个以上a结尾
		`(\d)\1{3}`,     // aaaa模式，包含4个以上a连续
	}
	skipFn := func(id uint64) (bool, error) {
		for _, p := range skipPattern {
			r := regexp2.MustCompile(p, 0) // 标准库regex不支持命名分组，所以第三方re库
			match, _ := r.MatchString(fmt.Sprintf("%d", id))
			if match {
				return true, nil
			}
		}
		return false, nil
	}

	locker := xlock.NewDLock("UidGenerator", cache.GetRedisClient(model.RedisDB))
	pool := cache2.NewUidQueuedPool("UidGenerator", cache.GetRedisClient(model.RedisDB))

	getMaxUid := func() (uint64, error) {
		id, err := dao.GetMaxUid()
		if err == nil && id < 1 {
			id = 100000
		}
		return id, err
	}

	var opts = []genuserid2.Option{
		genuserid2.WithSkipFunc(skipFn),
		//genuserid2.WithPoolConfig(10, 2),
	}
	uidGenerator = genuserid2.NewUidGenerator(locker, pool, getMaxUid, opts...)
	return nil
}
