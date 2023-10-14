package util

import (
	"context"
	"time"
)

var Ctx = context.Background()

func NewCtxWithTimeout(dur time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.TODO(), dur)
}
