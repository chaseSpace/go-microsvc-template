package handler

import (
	"github.com/hashicorp/go-uuid"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"microsvc/infra/xgrpc"
	"microsvc/pkg/xerr"
	"microsvc/pkg/xlog"
	"time"
)

type Handler func(ctx *fasthttp.RequestCtx) error
type UnaryInterceptor func(ctx *fasthttp.RequestCtx, handle Handler) error

func addInterceptor(handle func(ctx *fasthttp.RequestCtx) error, interceptor ...UnaryInterceptor) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		_ = interceptor[0](ctx, getChainUnaryHandler(interceptor, 0, handle))
	}
}

func getChainUnaryHandler(interceptors []UnaryInterceptor, curr int, finalInvoker Handler) Handler {
	if curr == len(interceptors)-1 {
		return finalInvoker
	}
	return func(ctx *fasthttp.RequestCtx) error {
		return interceptors[curr+1](ctx, getChainUnaryHandler(interceptors, curr+1, finalInvoker))
	}
}

// ------------ Interceptor ----------------

func logInterceptor(ctx *fasthttp.RequestCtx, handler Handler) (err error) {
	var (
		tid string
	)
	start := time.Now()
	defer func() {
		elapsed := time.Since(start).String()
		if xerr.IsNil(err) {
			xlog.Info("handle_ok", zap.ByteString("path", ctx.Path()), zap.String("dur", elapsed), zap.String("trace-id", tid))
		} else {
			xlog.Info("handle_fail", zap.ByteString("path", ctx.Path()), zap.Error(err), zap.String("dur", elapsed), zap.String("trace-id", tid))
		}
	}()
	err = handler(ctx)
	tid = ctx.UserValue(xgrpc.MetaKeyTraceId).(string)
	return err
}

func traceInterceptor(ctx *fasthttp.RequestCtx, handler Handler) (err error) {
	// TODO add trace logic
	tid, _ := uuid.GenerateUUID()
	ctx.SetUserValue(xgrpc.MetaKeyTraceId, tid)
	err = handler(ctx)
	return err
}
