package xlock

import (
	"context"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"microsvc/util"
	"time"
)

type DistributedLock interface {
	Lock(context.Context) error
	Unlock(context.Context) error
	IsLocked(context.Context) (bool, error)
}

var _ DistributedLock = (*DistributedLockInRedis)(nil)

type DistributedLockInRedis struct {
	client                  *redis.Client
	uniqueId, lockedRandVal string
}

const dLockKeyPrefix = "dLockKeyPrefix:"
const defaultLockExpiry = time.Minute

const unlockLuaScript = `
if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("del", KEYS[1])
else
    return 0
end
`

var (
	DLockFailed               = errors.New("redis lock failed")
	DLockFailedOnCtxTimeout   = errors.Wrap(DLockFailed, "ctx timeout")
	DLockFailedUpToMaxTimeout = errors.Wrap(DLockFailed, "up to max timeout")
	DUnlockFailed             = errors.New("redis unlock failed")
)

func NewDLock(uniqueId string, cli *redis.Client) DistributedLock {
	if len(uniqueId) < 5 {
		panic("uniqueId length must >= 5")
	}
	return &DistributedLockInRedis{
		client:   cli,
		uniqueId: dLockKeyPrefix + uniqueId,
	}
}

func (d *DistributedLockInRedis) Lock(ctx context.Context) (err error) {
	randVal := util.RandomString(6)

	// 设定一个最大超时
	var maxTimeoutCtx, cancel = context.WithTimeout(context.TODO(), time.Second)
	defer cancel()
	var cc = make(chan struct{})
	var lockOK bool
	for {

		go func() {
			defer func() { cc <- struct{}{} }()
			ret := d.client.SetNX(ctx, d.uniqueId, randVal, defaultLockExpiry)
			if ret.Err() != nil {
				err = errors.Wrap(DLockFailed, ret.Err().Error())
				return
			}
			if !ret.Val() {
				return
			}
			d.lockedRandVal = randVal
			lockOK = true
		}()

		select {
		case <-ctx.Done():
			return DLockFailedOnCtxTimeout
		case <-maxTimeoutCtx.Done():
			return DLockFailedUpToMaxTimeout
		case <-cc:
			if err != nil {
				return err
			}
			if lockOK {
				return
			}
		}
	}
}

func (d *DistributedLockInRedis) Unlock(ctx context.Context) error {
	if d.lockedRandVal == "" {
		return errors.Wrap(DUnlockFailed, "not locked")
	}
	ret := d.client.Eval(ctx, unlockLuaScript, []string{d.uniqueId}, d.lockedRandVal)
	i, err := ret.Int64()
	if err != nil {
		return errors.Wrap(DUnlockFailed, err.Error())
	}
	if i != 1 {
		return errors.Wrap(DUnlockFailed, "not locked or released")
	}
	return nil
}

func (d *DistributedLockInRedis) IsLocked(ctx context.Context) (bool, error) {
	if d.lockedRandVal == "" {
		return false, nil
	}
	ret := d.client.Exists(ctx, d.uniqueId)
	if ret.Err() != nil {
		return false, ret.Err()
	}
	return ret.Val() == 1, nil
}
