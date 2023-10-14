package genuserid

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

/*
递增 userid 生成模块（号池模式版本，支持高并发调用）
*/

var ErrPoolOperation = errors.New("pool operation func error")

type IncrementalPoolUIDGenerator struct {
	poolSizeThreshold int // 池id数量<=这个值就会扩容
	maxPoolSize       int // 池容量
	readyToPushIds    []uint64

	// 分布式锁（池扩容时使用）
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
	Pop() (uid uint64, size int, err error)
	Size() (size int, err error)
	MaxUnusedUID() (uid uint64, err error)
}

func (u *IncrementalPoolUIDGenerator) GenUid(ctx context.Context) (uint64, error) {
	var uid uint64
	var currPoolSize int
	var validId bool

	var cc = make(chan struct{})
	var err error

	//st := time.Now()
	//defer func() {
	//	println(333, time.Since(st).String())
	//}()
	// 最多只需要2次循环（第一次pool空则填充，第二次可从pool中获取到id）
	for i := 0; i < 2; i++ {
		go func() {
			defer func() {
				cc <- struct{}{}
			}()
			uid, currPoolSize, err = u.pool.Pop()
			if err != nil {
				return
			}
			//println(444, currPoolSize, time.Since(st).String())
			if uid == 0 || currPoolSize <= u.poolSizeThreshold {
				err = u.fillPool(currPoolSize)

				if err != nil {
					return
				}
				if uid == 0 {
					return // wait for next loop
				}
			}
			validId = true
		}()

		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-cc:
			if err != nil {
				return 0, err
			}
			if validId {
				return uid, nil
			}
		}
	}

	// 2次都拿不到id，且上面也没有err，说明是 NewUidGenerator 提供的池操作函数有问题，需要调用方自行检查
	return 0, ErrPoolOperation
}

// 扩容号池
// - 扩容时用了分布式锁，以提供并发取号性能
func (u *IncrementalPoolUIDGenerator) fillPool(currPoolSize int) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err = u.locker.Lock(ctx); err != nil {
		return errors.Wrap(err, "lock")
	}
	defer func() {
		err2 := u.locker.Unlock(ctx)
		if err2 != nil {
			err = errors.Wrap(err2, "unlock")
		}
	}()
	// 获得锁后，进行二次判断，避免并发时多次填充
	if currPoolSize, err = u.pool.Size(); err != nil {
		return errors.Wrapf(err, "pool size")
	} else if currPoolSize > u.poolSizeThreshold {
		return
	}

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
		u.readyToPushIds = u.readyToPushIds[:0] // reset
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
		u.readyToPushIds = append(u.readyToPushIds, maxUnusedUID)
	}
	//fmt.Printf("5555 %v %v\n", currPoolSize, u.readyToPushIds)
	return u.pool.Push(u.readyToPushIds)
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
