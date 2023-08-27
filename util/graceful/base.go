package graceful

import (
	"go.uber.org/zap"
	"microsvc/pkg/xlog"
	"microsvc/util"
	"os"
	"os/signal"
	"syscall"
)

var sigChan = make(chan os.Signal)

func SetupSignal() {
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
}

const logPrefix = "****** graceful ****** "

var stopFuncSlice []func()

func AddStopFunc(f func()) {
	stopFuncSlice = append(stopFuncSlice, f)
}

func OnExit() {
	stopAll() // case 2: backgroundSvc exited normally,  or signal received
	if err := recover(); err != nil {
		xlog.Panic(logPrefix+"server exited (thread panic)", zap.Any("err", err))
	}
	xlog.Info(logPrefix + "server exited")
	xlog.Stop()
}

func stopAll() {
	for _, stopF := range stopFuncSlice {
		stopF()
	}
}

func Run(backgroundSvc func()) {
	cc := util.Protect(backgroundSvc, func(err interface{}) {
		stopAll()                                                                  // case 1: backgroundSvc panic
		xlog.Panic(logPrefix+"backgroundSvc panic recovered", zap.Any("err", err)) // 再次panic
	})
	reason := ""
	select {
	case <-cc:
		reason = "(svc exited)"
	case <-sigChan:
		reason = "(signal)"
	}
	xlog.Warn(logPrefix + "server ready to exit" + reason)
}
