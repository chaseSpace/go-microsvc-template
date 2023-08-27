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
	select {
	case <-ctx.Done():
		return
	case <-quit:
		return
	}
}

func RunTaskWithTimeout(timeout time.Duration, f func()) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	RunTask(ctx, f)
}

func RunTaskWithCtxTimeout(timeout time.Duration, f func(ctx context.Context)) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	f(ctx)
}

func Protect(f func(), onPanic ...func(err interface{})) <-chan struct{} {
	exit := make(chan struct{})
	go func() {
		defer func() {
			if err := recover(); err != nil {
				if len(onPanic) > 0 {
					onPanic[0](err)
				}
			}
			exit <- struct{}{}
		}()
		f()
	}()
	return exit
}
