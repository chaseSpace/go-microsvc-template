package handler

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"microsvc/infra/xgrpc"
	"microsvc/pkg/xerr"
	"microsvc/pkg/xlog"
	"microsvc/pkg/xtime"
	"microsvc/protocol/svc"
	"microsvc/util"
	"time"
)

type Handler func(ctx *fasthttp.RequestCtx) ([]byte, error)
type UnaryInterceptor func(ctx *fasthttp.RequestCtx, handle Handler) ([]byte, error)

func addInterceptor(handle func(ctx *fasthttp.RequestCtx) ([]byte, error), interceptor ...UnaryInterceptor) fasthttp.RequestHandler {
	return func(fctx *fasthttp.RequestCtx) {
		fctx.SetContentType("application/json")
		fctx.SetStatusCode(200)
		fctx.Response.Header.Set("x-gateway-forward", "true")

		res, err := interceptor[0](fctx, getChainUnaryHandler(interceptor, 0, handle))
		if err == nil {
			fctx.SetBody(res) // transparent forwarding body
		} else {
			httpRes := &svc.GatewayHttpRsp{
				Code:        xerr.ErrInternal.Code,
				Msg:         xerr.ErrInternal.Msg,
				FromGateway: fctx.Value(ctxKeyFromGateway).(bool),
			}
			ecode, ok := xerr.FromErr(err)
			if ok {
				httpRes.Code = ecode.Code
				httpRes.Msg = ecode.Msg
			}
			fctx.SetBody(util.ToJson(httpRes))
		}
	}
}

func getChainUnaryHandler(interceptors []UnaryInterceptor, curr int, finalInvoker Handler) Handler {
	if curr == len(interceptors)-1 {
		return finalInvoker
	}
	return func(ctx *fasthttp.RequestCtx) ([]byte, error) {
		return interceptors[curr+1](ctx, getChainUnaryHandler(interceptors, curr+1, finalInvoker))
	}
}

// ------------ Interceptor ----------------

func logInterceptor(ctx *fasthttp.RequestCtx, handler Handler) (res []byte, err error) {
	tid := util.NewKsuid()
	ctx.SetUserValue(xgrpc.MetaKeyTraceId, tid)
	start := time.Now()
	xlog.Info("logInterceptor_start", zap.ByteString("path", ctx.Path()), zap.String("trace-id", tid))
	defer func() {

		elapsed := xtime.FormatDur(time.Since(start))
		if xerr.IsNil(err) {
			xlog.Info("handle_ok", zap.ByteString("path", ctx.Path()), zap.String("dur", elapsed), zap.String("trace-id", tid))
		} else {
			xlog.Info("handle_fail", zap.ByteString("path", ctx.Path()), zap.Error(err), zap.String("dur", elapsed), zap.String("trace-id", tid))
		}
	}()
	res, err = handler(ctx)
	return res, err
}

func traceInterceptor(ctx *fasthttp.RequestCtx, handler Handler) (res []byte, err error) {
	// TODO add trace logic
	//tid := ctx.UserValue(xgrpc.MetaKeyTraceId).(string)
	//println(111, tid)
	res, err = handler(ctx)
	return res, err
}
