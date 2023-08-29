package util

import (
	"context"
	"time"
)

func NewCtxWithTimeout(dur time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.TODO(), dur)
}
