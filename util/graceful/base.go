package graceful

import (
	"microsvc/deploy"
	"microsvc/pkg/xlog"
	"os"
	"os/signal"
	"syscall"
)

var sigChan = make(chan os.Signal, 1)

func Init() func(*deploy.XConfig, func(must bool, err error)) {
	return func(cc *deploy.XConfig, onEnd func(must bool, err error)) {
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		onEnd(true, nil)
	}
}

func OnExit(stop func()) {
	<-sigChan
	stop()
	xlog.Info("graceful: Server exited")
	os.Exit(0)
}
