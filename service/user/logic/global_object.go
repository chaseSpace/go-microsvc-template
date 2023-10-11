package logic

import (
	"fmt"
	"github.com/dlclark/regexp2"
	"microsvc/service/user/dao"
	"microsvc/util"
	"microsvc/util/db"
	"microsvc/xvendor/genuserid"
)

// 一些需要全局使用的资源在这里初始化
var (
	uidGenerator genuserid.UidGeneratorApi
)

func MustInit() {
	g := globalObjectCtrl{}
	util.AssertNilErr(g.InitUidGenerator())
}

type globalObjectCtrl struct {
}

func (globalObjectCtrl) InitUidGenerator() error {
	maxUID, err := dao.GetMaxUid()
	if db.IsMysqlErr(err) {
		return err
	}
	if maxUID < 1 {
		maxUID = 100000 // 6位数
	}
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

	uidGenerator = genuserid.NewUidGenerator(uint64(maxUID), dao.IsUidExists, skipFn)
	return nil
}
