package genuserid

import (
	"context"
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"sort"
	"sync"
	"testing"
	"time"
)

var _existIds []uint64
var _startId uint64

func resetCache() {
	_existIds = nil
	_startId = 0
}

func defaultGetCurrMaxUID() (uint64, error) {
	if len(_existIds) > 0 {
		return _existIds[len(_existIds)-1], nil
	}
	return _startId, nil
}

func timeoutCtx(to time.Duration) context.Context {
	c, _ := context.WithTimeout(context.Background(), to)
	return c
}

type Locker struct {
	mutex sync.Mutex
}

func (l *Locker) Lock(ctx context.Context) error {
	l.mutex.Lock()
	return nil
}
func (l *Locker) Unlock(ctx context.Context) error {
	l.mutex.Unlock()
	return nil
}

type queuedPool struct {
	idSlice []uint64
	lock    Locker // 池需要保证自己是并发安全的
}

func (q *queuedPool) MaxUnusedUID() (uid uint64, err error) {
	ctx := context.TODO()
	q.lock.Lock(ctx)
	defer q.lock.Unlock(ctx)
	if len(q.idSlice) > 0 {
		return q.idSlice[len(q.idSlice)-1], nil
	}
	return
}

func (q *queuedPool) Size() (size int, err error) {
	ctx := context.TODO()
	q.lock.Lock(ctx)
	defer q.lock.Unlock(ctx)
	return len(q.idSlice), nil
}

func (q *queuedPool) Push(ids []uint64) error {
	ctx := context.TODO()
	q.lock.Lock(ctx)
	defer q.lock.Unlock(ctx)
	q.idSlice = append(q.idSlice, ids...)
	return nil
}

func (q *queuedPool) Pop() (uint64, error) {
	ctx := context.TODO()
	q.lock.Lock(ctx)
	defer q.lock.Unlock(ctx)
	if len(q.idSlice) > 0 {
		id := q.idSlice[0]
		q.idSlice = q.idSlice[1:]
		return id, nil
	}
	return 0, nil
}

// SkipFn 是 NewUidGenerator 的可选参数之一，用来设置需要跳过的uid对应的规则
func TestSkipFn(t *testing.T) {
	// case-1. 不设置则完全自增
	skipFn := func(uint64) (bool, error) {
		return false, nil
	}
	g := NewUidGenerator(new(Locker), new(queuedPool), defaultGetCurrMaxUID, WithSkipFunc(skipFn))

	for i := 1; i <= 3; i++ {
		id, err := g.GenUid(context.TODO())
		assert.Nil(t, err)
		assert.Equal(t, uint64(i), id)

		// 生成一个可用id后，要存下来
		_existIds = append(_existIds, id)
	}

	// 清除缓存，以便不影响下一个测试
	resetCache()

	// case-2. 总是跳过，则一定会超时
	skipFn = func(uint64) (bool, error) {
		return true, nil
	}
	g = NewUidGenerator(new(Locker), new(queuedPool), defaultGetCurrMaxUID, WithSkipFunc(skipFn))

	for i := 0; i < 3; i++ {
		id, err := g.GenUid(timeoutCtx(time.Millisecond * 20))
		assert.Error(t, err, context.DeadlineExceeded)
		assert.Equal(t, uint64(0), id)
	}
	// 清除缓存，以便不影响下一个测试
	resetCache()

	// 3. 返回err，则GenUID会透传这个错误
	_err := errors.New("any")
	skipFn = func(id uint64) (bool, error) {
		return false, _err
	}

	g = NewUidGenerator(new(Locker), new(queuedPool), defaultGetCurrMaxUID, WithSkipFunc(skipFn))
	for i := 1; i <= 3; i++ {
		id, err := g.GenUid(context.TODO())
		assert.EqualError(t, errors.Wrap(_err, "skip err"), err.Error())
		assert.Equal(t, uint64(0), id)
	}
	// 清除缓存，以便不影响下一个测试
	resetCache()

	// 4. 正常逻辑
	var existIds, skipIds []uint64

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
	getCurrMaxUID := func() (uint64, error) {
		// 业务中此处应该从db读取
		if len(existIds) > 0 {
			return existIds[len(existIds)-1], nil
		}
		// 这里设置起始id
		return 9976, nil
	}
	g = NewUidGenerator(new(Locker), new(queuedPool), getCurrMaxUID, WithSkipFunc(skipFn))
	for i := 0; i < 23; i++ {
		id, err := g.GenUid(timeoutCtx(time.Millisecond * 20))
		assert.Nil(t, err)
		t.Log("Uid generated", id)

		existIds = append(existIds, id)
	}

	fmt.Printf("%+v\n", skipIds)

	// 每次进行池填充的都会把池填满，而默认池size=100，所以这里会跳过 (9976,10077] 范围内符合上述 skipFn 逻辑的id
	assert.Equal(t, []uint64{9977, 9988, 9999, 10000, 10011, 10022, 10033, 10044, 10055, 10066, 10077}, skipIds)

	resetCache()
}

type SafeIdBox struct {
	mutex sync.Mutex
	ids   []uint64
}

func (s *SafeIdBox) Add(id uint64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.ids = append(s.ids, id)
}

func (s *SafeIdBox) Last() (id uint64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(s.ids) > 0 {
		return s.ids[0]
	}
	return
}

// 并发测试
// 测试结果：   最低耗时  中位数  最高耗时
// 100个并发：  2ms     2.2ms   3.2ms
// 500个并发：  1.9ms   4.3ms   8.6ms
// 1000个并发： 2.5ms   9.3ms   16.8ms
func TestConcurrencyGenUID(t *testing.T) {
	var expectedIds, existIds, skipIds SafeIdBox
	expectedIdNum := 1000

	// 2. 设置要跳过的靓号模式（正则）,可选
	skipPattern := []string{
		`(\d)\1(\d)\2$`, // aabb结尾模式
		`(\d)\1{2}$`,    // aaa结尾模式，包含3个以上a结尾
		`(\d)\1{3}`,     // aaaa模式，包含4个以上a连续
	}
	skipFn := func(record bool) func(id uint64) (bool, error) {
		return func(id uint64) (bool, error) {
			for _, p := range skipPattern {
				r := regexp2.MustCompile(p, 0) // 标准库regex不支持分组的反向引用，所以第三方re库
				match, err := r.MatchString(fmt.Sprintf("%d", id))
				assert.Nil(t, err)
				if match {
					if record {
						skipIds.Add(id)
					}
					return true, nil
				}
			}
			return false, nil
		}
	}

	// 先把 expectedIds 填充到 expectedIdNum 个
	idStart := 1
	for i := 0; i < expectedIdNum; i++ {
		for {
			skip, err := skipFn(false)(uint64(idStart))
			if err != nil {
				t.Fatalf(err.Error())
			}
			if skip {
				idStart++
			} else {
				expectedIds.Add(uint64(idStart))
				break
			}
		}
		idStart++
	}

	assert.Equal(t, expectedIdNum, len(expectedIds.ids))

	// 这些都是可选的
	opts := []Option{
		WithSkipFunc(skipFn(true)),
		//WithPoolConfig(expectedIdNum, expectedIdNum/5), // 默认池大小是100，增加该值有助于提高并发
	}

	getCurrMaxUID := func() (uint64, error) {
		// 业务中此处应该从db读取
		if id := existIds.Last(); id > 0 {
			return id, nil
		}
		// 这里设置起始id
		return 0, nil
	}

	var x sync.WaitGroup
	var pool = new(queuedPool)
	var durationBox SafeIdBox
	g := NewUidGenerator(new(Locker), pool, getCurrMaxUID, opts...)
	for i := 0; i < expectedIdNum; i++ {
		x.Add(1)
		go func() {
			defer x.Done()
			st := time.Now()
			id, err := g.GenUid(timeoutCtx(time.Millisecond * 1000)) // 业务中建议设置1s
			durationBox.Add(uint64(time.Since(st)))
			assert.Nil(t, err)
			//t.Log("Uid generated", id)
			existIds.Add(id)
		}()
	}

	x.Wait()

	sort.SliceStable(skipIds.ids, func(i, j int) bool {
		return skipIds.ids[i] < skipIds.ids[j]
	})
	sort.SliceStable(existIds.ids, func(i, j int) bool {
		return existIds.ids[i] < existIds.ids[j]
	})
	sort.SliceStable(durationBox.ids, func(i, j int) bool {
		return durationBox.ids[i] < durationBox.ids[j]
	})

	// 这里的555，666，777，,888 是池扩充时跳过的id
	//assert.Equal(t, []uint64{111, 222, 333, 444, 555, 666, 777, 888}, skipIds.ids)
	assert.Equal(t, expectedIds.ids, existIds.ids)

	var durationStr []string
	for _, d := range durationBox.ids {
		durationStr = append(durationStr, time.Duration(d).String())
	}
	t.Logf("durations %+v", durationStr)
}

// 使用示例
func TestExample(t *testing.T) {
	var existIds, skipIds []uint64

	// 2. 设置要跳过的靓号模式（正则）,可选
	skipPattern := []string{
		`(\d)\1(\d)\2$`, // aabb结尾模式
		`(\d)\1{2}$`,    // aaa结尾模式，包含3个以上a结尾
		`(\d)\1{3}`,     // aaaa模式，包含4个以上a连续
	}
	skipFn := func(id uint64) (bool, error) {
		for _, p := range skipPattern {
			r := regexp2.MustCompile(p, 0) // 标准库regex不支持分组的反向引用，所以第三方re库
			match, err := r.MatchString(fmt.Sprintf("%d", id))
			assert.Nil(t, err)
			if match {
				skipIds = append(skipIds, id)
				return true, nil
			}
		}
		return false, nil
	}

	// 这些都是可选的
	opts := []Option{
		// 指定需要跳过的UID（比如靓号逻辑）
		WithSkipFunc(skipFn),
	}

	getCurrMaxUID := func() (uint64, error) {
		// 业务中此处应该从db读取
		if len(existIds) > 0 {
			return existIds[len(existIds)-1], nil
		}
		// 这里设置起始id
		return 0, nil
	}

	g := NewUidGenerator(new(Locker), new(queuedPool), getCurrMaxUID, opts...)
	for i := 0; i < 10; i++ {
		id, err := g.GenUid(timeoutCtx(time.Second)) // 超时建议1s，主要考虑池操作本身是db操作，通常会稍微耗时
		assert.Nil(t, err)
		t.Log("Uid generated", id)

		existIds = append(existIds, id)
	}

	assert.Equal(t, []uint64(nil), skipIds)
	assert.Equal(t, []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, existIds)
}
