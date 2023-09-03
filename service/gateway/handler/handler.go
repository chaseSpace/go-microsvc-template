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
	"microsvc/protocol/svc"
	"microsvc/util"
	"net/http"
	"regexp"
	"time"
)

type GatewayCtrl struct {
}

const applicationJson = "application/json"

const gatewayForwardTimeout = time.Second * 5

func (GatewayCtrl) Handler(ctx *fasthttp.RequestCtx) {
	interceptors := []UnaryInterceptor{logInterceptor, traceInterceptor}
	addInterceptor(forwardHandler, interceptors...)(ctx)
}

var (
	routerRegexToSvc = regexp.MustCompile(`svc.(\w+).(\w+)Ext`)
)

func forwardHandler(fctx *fasthttp.RequestCtx) error {

	fullPath := string(fctx.Path())
	items := routerRegexToSvc.FindStringSubmatch(fullPath)
	if len(items) != 3 {
		return xerr.ErrApiNotFound
	}

	var (
		service = enums.Svc(items[1])
		errcode = xerr.ErrNil
	)

	var (
		// TODO: optimize with pool
		res = bytes.NewBuffer(make([]byte, 0, 512))
		// fromGateway 表示是否直接从网关响应，而不进行进一步的转发
		fromGateWay = true
	)

	defer func() {
		fctx.SetContentType(applicationJson)
		fctx.SetStatusCode(http.StatusOK)

		if errcode.IsNil() {
			fctx.SetBody(res.Bytes()) // transparent forwarding body
		} else {
			httpRes := &svc.GatewayHttpRsp{Code: errcode.Code, Msg: errcode.Msg, FromGateway: fromGateWay}
			fctx.SetBody(util.ToJson(httpRes))
		}
	}()

	conn := svccli.GetConn(service)
	if conn == nil {
		errcode = xerr.ErrNoRPCClient.AppendMsg(service.Name())
		return errcode
	}

	// below is grpc calling, we set `fromGateWay` to false whether the call returns an error or not
	fromGateWay = false

	ctx, cancel := newRpcCtx(fctx)
	defer cancel()

	err := conn.Invoke(ctx, fullPath, fctx.PostBody(), res, grpc.CallContentSubtype(protobytes.Bytes))
	if err != nil {
		errcode = err.(xerr.XErr) // err is converted to XErr in grpc client interceptor
	}
	return errcode
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
