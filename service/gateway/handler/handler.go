package handler

import (
	"bytes"
	"context"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"microsvc/bizcomm/auth"
	"microsvc/enums/svc"
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

var (
	interceptors = []UnaryInterceptor{interceptor{}.Trace, interceptor{}.Log}
)

func (GatewayCtrl) Handler(ctx *fasthttp.RequestCtx) {
	addInterceptor(forwardHandler, interceptors...)(ctx)
}

// ----------------------------------------------------------------

const (
	pathProxyPrefix   = "/forward/"
	pathPing          = "/ping"
	ctxKeyFromGateway = "from-gateway"
)

var (
	routeRegexToSvc = regexp.MustCompile(`svc\.(\w+)\.\w+Ext/\w+`)
	//routeRegexToAdmin = regexp.MustCompile(`admin\.(\w+)\.\w+Ext/\w+`)
)

func forwardHandler(fctx *fasthttp.RequestCtx) ([]byte, error) {
	fctx.SetUserValue(ctxKeyFromGateway, true)
	var (
		// TODO: optimize with pool
		res = bytes.NewBuffer(make([]byte, 0, 512))
	)

	fullPath := string(fctx.Path())
	if !strings.HasPrefix(fullPath, pathProxyPrefix) {
		switch fullPath {
		case pathPing: // for test
			return []byte("pong"), nil
		}
		return nil, xerr.ErrNotFound.New("path must start with %s", pathProxyPrefix)
	}

	dstPath := fullPath[len(pathProxyPrefix):]
	items := routeRegexToSvc.FindStringSubmatch(dstPath)
	if len(items) != 2 {
		return nil, xerr.ErrNotFound
	}

	fctx.SetUserValue(ctxKeyFromGateway, false)
	var (
		service     = svc.Svc(items[1])
		forwardPath = items[0]
	)

	conn := svccli.GetConn(service)

	ctx, cancel := newRpcCtx(fctx)
	defer cancel()

	err := conn.Invoke(ctx, forwardPath, fctx.PostBody(), res, grpc.CallContentSubtype(protobytes.Name))
	if err != nil {
		return nil, err.(xerr.XErr) // err is converted to XErr in grpc client interceptor
	}
	return res.Bytes(), nil
}

func newRpcCtx(fctx *fasthttp.RequestCtx) (context.Context, context.CancelFunc) {

	traceId, _ := fctx.Value(xgrpc.MdKeyTraceId).(string)

	md := metadata.Pairs(
		xgrpc.MdKeyAuthToken, string(fctx.Request.Header.Peek(auth.HeaderKey)),
		xgrpc.MdKeyTraceId, traceId,
		xgrpc.MdKeyFromGatewayFlag, xgrpc.MdKeyFlagExist,
	)

	ctx, cancel := context.WithTimeout(context.TODO(), gatewayForwardTimeout)

	return metadata.NewOutgoingContext(ctx, md), cancel
}
