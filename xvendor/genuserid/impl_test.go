package genuserid

import (
	"context"
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func defaultExistFn(uint64) (bool, error) {
	return false, nil
}

func defaultSkipFn(uint64) (bool, error) {
	return false, nil
}

func timeoutCtx(to time.Duration) context.Context {
	c, _ := context.WithTimeout(context.Background(), to)
	return c
}

func TestExistFn(t *testing.T) {
	// 1.always false, then always get same startUID
	existFn := func(uint64) (bool, error) {
		return false, nil
	}
	g := NewUidGenerator(1, existFn, defaultSkipFn)

	for i := 0; i < 3; i++ {
		id, err := g.GenUid(context.TODO())
		assert.Nil(t, err)
		assert.Equal(t, uint64(1), id)
		g.UpdateStartUid(id)
	}

	// 2.always true, then timeout
	existFn = func(uint64) (bool, error) {
		return true, nil
	}
	g = NewUidGenerator(1, existFn, defaultSkipFn)

	for i := 0; i < 3; i++ {
		id, err := g.GenUid(timeoutCtx(time.Millisecond * 20))
		assert.ErrorIs(t, err, context.DeadlineExceeded)
		assert.Equal(t, uint64(0), id)
	}

	// 3. existFn returns err
	_err := errors.New("any")
	existFn = func(id uint64) (bool, error) {
		return false, _err
	}

	g = NewUidGenerator(1, existFn, defaultSkipFn)
	for i := 1; i <= 3; i++ {
		id, err := g.GenUid(context.TODO())
		assert.EqualError(t, err, _err.Error())
		assert.Equal(t, uint64(0), id)
	}

	// 4.normal logic
	var uids []uint64
	existFn = func(id uint64) (bool, error) {
		return lo.Contains(uids, id), nil
	}

	g = NewUidGenerator(1, existFn, defaultSkipFn)
	for i := 1; i <= 3; i++ {
		id, err := g.GenUid(context.TODO())
		assert.Nil(t, err)
		assert.Equal(t, uint64(i), id)
		g.UpdateStartUid(id)

		uids = append(uids, id)
	}
}

func TestSkipFn(t *testing.T) {
	// 1.always false, then always get same startUID
	skipFn := func(uint64) (bool, error) {
		return false, nil
	}
	g := NewUidGenerator(1, defaultExistFn, skipFn)

	for i := 0; i < 3; i++ {
		id, err := g.GenUid(context.TODO())
		assert.Nil(t, err)
		assert.Equal(t, uint64(1), id)
		g.UpdateStartUid(id)
	}

	// 2.always true, then timeout
	skipFn = func(uint64) (bool, error) {
		return true, nil
	}
	g = NewUidGenerator(1, defaultExistFn, skipFn)

	for i := 0; i < 3; i++ {
		id, err := g.GenUid(timeoutCtx(time.Millisecond * 20))
		assert.ErrorIs(t, err, context.DeadlineExceeded)
		assert.Equal(t, uint64(0), id)
	}

	// 3. skipFn returns err
	_err := errors.New("any")
	skipFn = func(id uint64) (bool, error) {
		return false, _err
	}

	g = NewUidGenerator(1, defaultExistFn, skipFn)
	for i := 1; i <= 3; i++ {
		id, err := g.GenUid(context.TODO())
		assert.EqualError(t, err, _err.Error())
		assert.Equal(t, uint64(0), id)
	}

	// 4.normal logic
	var existIds, skipIds []uint64
	existFn := func(id uint64) (bool, error) {
		return lo.Contains(existIds, id), nil
	}

	skipPattern := []string{
		`(\d)\1(\d)\2$`, // aabb结尾模式
		`(\d)\1{2}$`,    // aaa结尾模式，包含3个以上a结尾
		`(\d)\1{3}`,     // aaaa模式，包含4个以上a连续
	}
	skipFn = func(id uint64) (bool, error) {
		for _, p := range skipPattern {
			r := regexp2.MustCompile(p, 0) // 标准库regex不支持分组引用，所以使用第三方re库
			match, err := r.MatchString(fmt.Sprintf("%d", id))
			assert.Nil(t, err)
			if match {
				skipIds = append(skipIds, id)
				return true, nil
			}
		}
		return false, nil
	}

	g = NewUidGenerator(9977, existFn, skipFn, WithLimitOnceLoopTimes(10))
	for i := 0; i < 23; i++ {
		id, err := g.GenUid(timeoutCtx(time.Millisecond * 20))
		assert.Nil(t, err)
		t.Log("Uid generated", id)
		g.UpdateStartUid(id)

		existIds = append(existIds, id)
	}
	assert.Equal(t, []uint64{9977, 9988, 9999, 10000}, skipIds)
}

// 使用示例
func TestExample(t *testing.T) {
	var existIds, skipIds []uint64
	// 1. 定义用来判断已经用过的uid的函数
	existFn := func(id uint64) (bool, error) {
		// 一般通过尝试写入db是否成功 来判断是否存在
		// 省略。。。
		return lo.Contains(existIds, id), nil
	}

	// 2. 设置要跳过的靓号模式（正则）
	skipPattern := []string{
		`(\d)\1(\d)\2$`, // aabb结尾模式
		`(\d)\1{2}$`,    // aaa结尾模式，包含3个以上a结尾
		`(\d)\1{3}`,     // aaaa模式，包含4个以上a连续
	}
	skipFn := func(id uint64) (bool, error) {
		for _, p := range skipPattern {
			r := regexp2.MustCompile(p, 0) // 标准库regex不支持命名分组，所以第三方re库
			match, err := r.MatchString(fmt.Sprintf("%d", id))
			assert.Nil(t, err)
			if match {
				skipIds = append(skipIds, id)
				return true, nil
			}
		}
		return false, nil
	}

	// 3. 初始化，设置起始id，比如是6位数id，则设置100000（仅首次），正常应该从db中读取当前最大已使用id+1
	readFromDB := func() uint64 {
		return 100000
	}
	startUid := readFromDB()
	g := NewUidGenerator(startUid, existFn, skipFn)
	for i := 0; i < 10; i++ {
		id, err := g.GenUid(timeoutCtx(time.Second * 2)) // 超时建议2s，主要考虑existFn通常会读取数据较为耗时
		assert.Nil(t, err)
		t.Log("Uid generated", id)
		g.UpdateStartUid(id)

		existIds = append(existIds, id)
	}
	t.Log("Uid skipped", skipIds)
}
