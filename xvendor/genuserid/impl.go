package genuserid

import (
	"context"
	"errors"
	"sync"
)

/*
userid 生成模块
*/

var ErrReachOnceLoopTimesLimit = errors.New("reach to max once loop times limit")

type UidGenerator struct {
	startUid           uint64
	limitOnceLoopTimes int
	lock               sync.Mutex
	existFn            func(uint64) (bool, error)
	skipFn             func(uint64) (bool, error)
}

func (u *UidGenerator) UpdateStartUid(id uint64) {
	u.startUid = id
	return
}

func (u *UidGenerator) GenUid(ctx context.Context) (uint64, error) {
	u.lock.Lock()
	defer u.lock.Unlock()
	id := u.startUid

	validId := false
	var cc = make(chan struct{})
	var boolVal bool
	var err error
	for i := 1; ; i++ {
		go func(id uint64) {
			defer func() {
				cc <- struct{}{}
			}()
			if boolVal, err = u.existFn(id); err != nil || boolVal {
				return
			}
			if boolVal, err = u.skipFn(id); err != nil || boolVal {
				return
			}
			validId = true
		}(id)

		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-cc:
			if err != nil {
				return 0, err
			}
			if validId {
				return id, nil
			}
			if u.limitOnceLoopTimes > 0 && u.limitOnceLoopTimes == i {
				return id, ErrReachOnceLoopTimesLimit
			}
			id++
		}
	}
}

type Option func(generator *UidGenerator)

func WithLimitOnceLoopTimes(i int) Option {
	return func(generator *UidGenerator) {
		generator.limitOnceLoopTimes = i
	}
}

func NewUidGenerator(startUid uint64, existFn, skipFn func(uint64) (bool, error), opts ...Option) UidGeneratorApi {
	if startUid == 0 {
		panic("startUid must be non-zero")
	}
	g := &UidGenerator{startUid: startUid}
	for _, opt := range opts {
		opt(g)
	}
	g.existFn = existFn
	g.skipFn = skipFn
	return g
}
