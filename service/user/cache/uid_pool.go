package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"microsvc/infra/cache"
	"microsvc/util"
	"microsvc/xvendor/genuserid2"
)

type UidQueuedPool struct {
	redisKey string
	client   *redis.Client
}

var _ genuserid2.QueuedPool = (*UidQueuedPool)(nil)

func NewUidQueuedPool(key string, client *redis.Client) *UidQueuedPool {
	return &UidQueuedPool{
		redisKey: key,
		client:   client,
	}
}

func (u *UidQueuedPool) Push(ids []uint64) error {
	ids2 := lo.Map(ids, func(item uint64, index int) interface{} {
		return item
	})
	return u.client.LPush(context.TODO(), u.redisKey, ids2).Err()
}

func (u *UidQueuedPool) Pop() (uid uint64, err error) {
	ret := u.client.RPop(util.Ctx, u.redisKey)
	ignored, err := cache.IgnoreNil(ret.Err())
	if err != nil {
		return 0, err
	}
	if ignored {
		return 0, nil
	}
	uid2, err := ret.Uint64()
	if err != nil {
		return 0, err
	}
	return uid2, nil
}

func (u *UidQueuedPool) Size() (size int, err error) {
	ret := u.client.LLen(util.Ctx, u.redisKey)
	if ret.Err() != nil {
		return 0, ret.Err()
	}
	return int(ret.Val()), nil
}

func (u *UidQueuedPool) MaxUnusedUID() (uid uint64, err error) {
	ret := u.client.LIndex(util.Ctx, u.redisKey, 0)
	if ret.Err() != nil {
		_, err = cache.IgnoreNil(ret.Err())
		return
	}
	return ret.Uint64()
}

func (u *UidQueuedPool) Reset() error {
	return u.client.Del(util.Ctx, u.redisKey).Err()
}
