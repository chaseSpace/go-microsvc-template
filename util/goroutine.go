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
