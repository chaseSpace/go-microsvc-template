package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"testing"
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

func TestUidQueuedPool_PopPush(t *testing.T) {
	client := initRedis()
	defer client.Close()

	p := NewUidQueuedPool("uid_pool", client)
	_ = p.Reset()

	// pop without no push
	uid, err := p.Pop()
	assert.Nil(t, err)
	assert.Zero(t, uid)

	// push 1,2,3, pop 1,2,3
	ids := []uint64{1, 2, 3}
	err = p.Push(ids)
	assert.Nil(t, err)

	lo.ForEach(ids, func(item uint64, index int) {
		uid, err = p.Pop()
		assert.Nil(t, err)
		assert.EqualValues(t, item, uid)
	})

	// push 4,5,6, pop 4,5, then push 100
	p.Push([]uint64{4, 5, 6})
	p.Pop()
	p.Pop()
	p.Push([]uint64{100})

	// pop return 6,100
	uid, _ = p.Pop()
	assert.EqualValues(t, 6, uid)
	uid, _ = p.Pop()
	assert.EqualValues(t, 100, uid)
}

func TestUidQueuedPool_MaxUnusedUID(t *testing.T) {
	client := initRedis()
	defer client.Close()

	p := NewUidQueuedPool("uid_pool", client)
	_ = p.Reset()

	// no data
	maxId, err := p.MaxUnusedUID()
	assert.Nil(t, err)
	assert.EqualValues(t, 0, maxId)

	// push 1,2,3, pop 1,2,3
	ids := []uint64{1, 2, 3}
	err = p.Push(ids)
	assert.Nil(t, err)

	maxId, err = p.MaxUnusedUID()
	assert.Nil(t, err)
	assert.EqualValues(t, 3, maxId)

	// pop 3, push 100, 101
	p.Pop()
	p.Push([]uint64{100, 101})

	maxId, err = p.MaxUnusedUID()
	assert.Nil(t, err)
	assert.EqualValues(t, 101, maxId)
}

func TestUidQueuedPool_Size(t *testing.T) {
	client := initRedis()
	defer client.Close()

	p := NewUidQueuedPool("uid_pool", client)
	_ = p.Reset()

	size, err := p.Size()
	assert.Nil(t, err)
	assert.EqualValues(t, 0, size)

	// push 1,2,3, got size equal to 3
	ids := []uint64{1, 2, 3}
	err = p.Push(ids)
	assert.Nil(t, err)

	size, err = p.Size()
	assert.Nil(t, err)
	assert.EqualValues(t, len(ids), size)

	// pop 1, got size equal to 2
	p.Pop()
	size, err = p.Size()
	assert.Nil(t, err)
	assert.EqualValues(t, len(ids)-1, size)

	// pop 2,3, size -> 0
	p.Pop()
	p.Pop()
	size, err = p.Size()
	assert.Nil(t, err)
	assert.EqualValues(t, 0, size)
}
