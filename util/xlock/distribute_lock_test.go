package xlock

import (
	"context"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func initRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:       "127.0.0.1:6379",
		Password:   "123",
		DB:         0,
		MaxRetries: 2,
	})
	err := rdb.Ping(context.TODO()).Err()
	if err != nil {
		panic(err)
	}
	return rdb
}

func timeoutCtx(to time.Duration) context.Context {
	c, _ := context.WithTimeout(context.Background(), to)
	return c
}

func TestNewDLockBasic(t *testing.T) {
	cli := initRedis()
	defer cli.Close()

	k := NewDLock("hello_lock", cli)

	var imap = make(map[int]bool)

	count := 100 // 设置过高 会返回超时错误
	v := 0

	var x sync.WaitGroup
	for i := 0; i < count; i++ {

		x.Add(1)
		go func() {
			defer x.Done()
			err := k.Lock(context.TODO())
			if !assert.Nil(t, err) {
				t.FailNow()
			}

			v++
			imap[v] = true
			err = k.Unlock(context.TODO())
			if !assert.Nil(t, err) {
				t.FailNow()
			}
		}()

	}

	x.Wait()
	assert.Equal(t, count, len(imap))
}

func TestNewDLockWithSetCtxTimeout(t *testing.T) {
	cli := initRedis()
	defer cli.Close()

	k := NewDLock("hello_lock", cli)

	// case-1:lock timeout
	err := k.Lock(timeoutCtx(time.Microsecond)) // 1微秒 必定超时
	if !assert.Equal(t, DLockFailedOnCtxTimeout, err) {
		t.FailNow()
	}

	err = k.Unlock(context.TODO())
	if !assert.EqualError(t, errors.Wrap(DUnlockFailed, "not locked"), err.Error()) {
		t.FailNow()
	}
}

func TestNewDLockWithLockAgainAndNotUnlock(t *testing.T) {
	cli := initRedis()
	defer cli.Close()

	k := NewDLock("hello_lock", cli)

	err := k.Lock(context.TODO())
	assert.Nil(t, err)

	// try to lock again
	err = k.Lock(context.TODO())
	assert.Equal(t, DLockFailedUpToMaxTimeout, err)

	err = k.Unlock(context.TODO())
	assert.Nil(t, err)

	// lock successfully after unlock
	err = k.Lock(context.TODO())
	assert.Nil(t, err)
	err = k.Unlock(context.TODO())
	assert.Nil(t, err)
}

func TestNewDLockWithUnlockWithoutLock(t *testing.T) {
	cli := initRedis()
	defer cli.Close()

	k := NewDLock("hello_lock", cli)

	err := k.Unlock(context.TODO())
	assert.EqualError(t, errors.Wrap(DUnlockFailed, "not locked"), err.Error())
}

func TestNewDLockWithIsLocked(t *testing.T) {
	cli := initRedis()
	defer cli.Close()

	k := NewDLock("hello_lock", cli)

	// IsLocked returns false if not locked
	exists, err := k.IsLocked(context.TODO())
	assert.Nil(t, err)
	assert.False(t, exists)

	// then we get lock
	err = k.Lock(context.TODO())
	assert.Nil(t, err)

	// now, IsLocked returns true
	exists, err = k.IsLocked(context.TODO())
	assert.Nil(t, err)
	assert.False(t, !exists)
}
