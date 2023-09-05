package handler

import (
	"bytes"
	"context"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"microsvc/enums"
	"microsvc/infra/svccli"
	"microsvc/infra/xgrpc"
	"microsvc/infra/xgrpc/protobytes"
	"microsvc/pkg/xerr"
	"regexp"
	"strings"
	"time"
)

type GatewayCtrl struct {
}

const gatewayForwardTimeout = time.Second * 5

func (GatewayCtrl) Handler(ctx *fasthttp.RequestCtx) {
	interceptors := []UnaryInterceptor{logInterceptor, traceInterceptor}
	addInterceptor(forwardHandler, interceptors...)(ctx)
}

// ----------------------------------------------------------------

const (
	apiUnionPathPrefix = "/forward/"
	ctxKeyFromGateway  = "from-gateway"
)

var (
	routerRegexToSvc = regexp.MustCompile(`svc.(\w+).(\w+)Ext/\w+`)
)

func forwardHandler(fctx *fasthttp.RequestCtx) ([]byte, error) {
	fctx.SetUserValue(ctxKeyFromGateway, true)
	var (
		// TODO: optimize with pool
		res = bytes.NewBuffer(make([]byte, 0, 512))
	)

	fullPath := string(fctx.Path())
	if !strings.HasPrefix(fullPath, apiUnionPathPrefix) {
		return nil, xerr.ErrApiNotFound.NewMsg("path must start with %s", apiUnionPathPrefix)
	}

	dstPath := fullPath[len(apiUnionPathPrefix):]
	items := routerRegexToSvc.FindStringSubmatch(dstPath)
	if len(items) != 3 {
		return nil, xerr.ErrApiNotFound
	}

	fctx.SetUserValue(ctxKeyFromGateway, false)
	var (
		service     = enums.Svc(items[2])
		forwardPath = items[0]
	)
	conn := svccli.GetConn(service)
	if conn == nil {
		return nil, xerr.ErrNoRPCClient.AppendMsg(service.Name())
	}

	ctx, cancel := newRpcCtx(fctx)
	defer cancel()

	err := conn.Invoke(ctx, forwardPath, fctx.PostBody(), res, grpc.CallContentSubtype(protobytes.Bytes))
	if err != nil {
		return nil, err.(xerr.XErr) // err is converted to XErr in grpc client interceptor
	}
	return res.Bytes(), nil
}

func newRpcCtx(fctx *fasthttp.RequestCtx) (context.Context, context.CancelFunc) {

	traceId, _ := fctx.Value(xgrpc.MetaKeyTraceId).(string)

	md := metadata.Pairs(
		xgrpc.MetaKeyFromGateway, "true",
		xgrpc.MetaKeyAuth, "Bearer 123",
		xgrpc.MetaKeyTraceId, traceId,
	)

	ctx, cancel := context.WithTimeout(context.TODO(), gatewayForwardTimeout)

	return metadata.NewOutgoingContext(ctx, md), cancel
}
