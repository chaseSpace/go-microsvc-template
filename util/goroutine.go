package util

import (
	"context"
	"time"
)

func RunTask(ctx context.Context, f func()) {
	quit := make(chan struct{})
	go func() {
		f()
		quit <- struct{}{}
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case <-quit:
			return
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}
