package genuserid2

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

/*
递增 userid 生成模块（号池版本，支持高并发调用）
*/

type IncrementalPoolUIDGenerator struct {
	poolSizeThreshold int // 池剩余id数量<=这个值就会扩容
	maxPoolSize       int // 池容量
	pushIdsBuffer     []uint64

	// 分布式锁
	locker DistributeLock

	getCurrentMaxUID func() (uint64, error)
	skipFn           func(uint64) (bool, error)
	pool             QueuedPool
}

type DistributeLock interface {
	Lock(ctx context.Context) error
	Unlock(ctx context.Context) error
}

type QueuedPool interface {
	Push(ids []uint64) error
	Pop() (uid uint64, err error)
	Size() (size int, err error)
	MaxUnusedUID() (uid uint64, err error)
}

func (u *IncrementalPoolUIDGenerator) GenUid(ctx context.Context) (uid uint64, err error) {
	var cc = make(chan struct{})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err = u.locker.Lock(ctx); err != nil {
		return 0, errors.Wrap(err, "lock")
	}
	defer func() {
		err2 := u.locker.Unlock(ctx)
		if err2 != nil {
			err = errors.Wrap(err2, "unlock")
		}
	}()

	var size int
	//st := time.Now()
	//defer func() {
	//	println(333, time.Since(st).String())
	//}()

	// 当池容量 远小于 并发请求数时，循环次数会>2（所以要根据预估的业务的并发请求数 来配置 池容量，避免此操作阻塞过久）
	for i := 0; ; i++ {
		go func() {
			defer func() {
				cc <- struct{}{}
			}()
			uid, err = u.pool.Pop()
			if err != nil {
				return
			}

			size, err = u.pool.Size()
			if err != nil {
				return
			}
			//println(444, i, uid, size)

			if uid == 0 || size <= u.poolSizeThreshold {
				err = u.fillPool(size)
			}
		}()

		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-cc:
			if err != nil {
				return
			}
			if uid > 0 {
				return
			}
		}
	}
}

// 填充号池
// - 填充时用了分布式锁，以保证并发安全
func (u *IncrementalPoolUIDGenerator) fillPool(currPoolSize int) (err error) {
	maxUnusedUID := uint64(0)
	if currPoolSize > 0 {
		maxUnusedUID, err = u.pool.MaxUnusedUID()
		if err != nil {
			return errors.Wrapf(err, "MaxUnusedUID")
		}
	}

	if maxUnusedUID == 0 {
		// 获取业务中当前已使用的最大uid，进行对比，取较大值
		if currMaxUID, err := u.getCurrentMaxUID(); err != nil {
			return errors.Wrap(err, "getCurrentMaxUID")
		} else {
			maxUnusedUID = currMaxUID
		}
	}

	defer func() {
		u.pushIdsBuffer = u.pushIdsBuffer[:0] // reset
	}()

	// 确保每次把号池填满
	for i := 0; i < u.maxPoolSize-currPoolSize; i++ {
		for {
			maxUnusedUID++
			if skip, err := u.skipFn(maxUnusedUID); err != nil {
				return errors.Wrap(err, "skip err")
			} else if skip {
				continue
			}
			break
		}
		u.pushIdsBuffer = append(u.pushIdsBuffer, maxUnusedUID)
	}
	//fmt.Printf("5555 %v %v\n",  currPoolSize, u.pushIdsBuffer)
	return u.pool.Push(u.pushIdsBuffer)
}

type Option func(generator *IncrementalPoolUIDGenerator)

func WithPoolConfig(maxPoolSize, thresholdSize int) Option {
	return func(generator *IncrementalPoolUIDGenerator) {
		generator.maxPoolSize = maxPoolSize
		generator.poolSizeThreshold = thresholdSize
	}
}

func WithSkipFunc(skipFn func(uint64) (bool, error)) Option {
	return func(generator *IncrementalPoolUIDGenerator) {
		generator.skipFn = skipFn
	}
}

func NewUidGenerator(locker DistributeLock, pool QueuedPool, getCurrMaxUID func() (uint64, error), opts ...Option) UIDGeneratorApi {
	g := &IncrementalPoolUIDGenerator{locker: locker, pool: pool, getCurrentMaxUID: getCurrMaxUID}
	for _, opt := range opts {
		opt(g)
	}
	if g.maxPoolSize < 10 {
		g.maxPoolSize = 100
	}
	if g.poolSizeThreshold < 1 {
		g.poolSizeThreshold = g.maxPoolSize / 5
	}
	return g
}
