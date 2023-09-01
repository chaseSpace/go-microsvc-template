package handler

import (
	"bytes"
	"context"
	"fmt"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"microsvc/enums"
	"microsvc/infra/svccli"
	"microsvc/infra/xgrpc"
	"microsvc/infra/xgrpc/proto"
	"microsvc/pkg/xerr"
	"microsvc/protocol/svc"
	"microsvc/util"
	"net/http"
	"strings"
	"sync"
	"time"
)

type GatewayCtrl struct {
}

const applicationJson = "application/json"

const gatewayForwardTimeout = time.Second * 5

func (GatewayCtrl) Handler(ctx *fasthttp.RequestCtx) {
	addInterceptor(forwardHandler, logInterceptor, traceInterceptor)(ctx)
}

func forwardHandler(fctx *fasthttp.RequestCtx) error {
	errtyp := xerr.ErrOK

	var (
		v = bodyPool.Get().(*reqAndRes)
		// TODO: optimize within pool
		res = bytes.NewBuffer(make([]byte, 0, 512))
		// if true,that indicates respond from gateway directly, not forwarding yet.
		fromGateWay = true
	)
	path := string(fctx.Path())
	defer func() {
		bodyPool.Put(v) // return to pool

		fctx.SetContentType(applicationJson)
		fctx.SetStatusCode(http.StatusOK)

		if errtyp.IsOK() {
			fctx.SetBody(res.Bytes()) // transparent forwarding body
		} else {
			httpRes := &svc.GatewayHttpRsp{Code: errtyp.Code, Msg: errtyp.Msg, FromGateway: fromGateWay}
			fctx.SetBody(util.ToJson(httpRes))
		}
	}()

	route := parseRoute(strings.TrimLeft(path, "/"))
	if route == nil {
		errtyp = xerr.ErrApiNotFound
		return errtyp
	}

	conn := svccli.GetConn(route.Svc)
	if conn == nil {
		errtyp = xerr.ErrNoRPCClient.AppendMsg(route.Svc.Name())
		return errtyp
	}

	// below is grpc calling, we set `fromGateWay` to false whether the call returns an error or not
	fromGateWay = false
	var (
		traceId, _  = fctx.Value(xgrpc.MetaKeyTraceId).(string)
		ctx, cancel = context.WithTimeout(context.TODO(), gatewayForwardTimeout)
	)
	defer cancel()
	err := conn.Invoke(newRpcCtx(ctx, traceId), path, fctx.PostBody(), res, grpc.CallContentSubtype(proto.Bytes))
	if err != nil {
		// err is converted to XErr in grpc client interceptor
		errtyp = err.(xerr.XErr)
	}
	return errtyp
}

type SvcApiRoute struct {
	Svc    enums.Svc
	Prefix string
	Method string
}

func (r SvcApiRoute) UnionMethod() string {
	return fmt.Sprintf("%s/Forward", r.Prefix)
}

// e.g. path is "svc.user.UserExt/GetUser"
func parseRoute(path string) *SvcApiRoute {
	if !strings.HasPrefix(path, "svc.") {
		return nil
	}
	ss := strings.Split(path, "/")
	if len(ss) == 2 && strings.HasSuffix(ss[0], "Ext") {
		ss2 := strings.Split(ss[0], ".")
		if len(ss2) == 3 && len(ss2[1]) <= 20 {
			return &SvcApiRoute{Svc: enums.Svc(ss2[1]), Prefix: ss[0], Method: ss[1]}
		}
	}
	return nil
}

type reqAndRes struct {
	Req, Res proto.ArbitraryBody
}

var bodyPool = sync.Pool{
	New: func() interface{} {
		return &reqAndRes{}
	},
}

func newRpcCtx(ctx context.Context, traceId string) context.Context {
	md := metadata.Pairs(
		xgrpc.MetaKeyFromGateway, "true",
		xgrpc.MetaKeyAuth, "Bearer 123",
		xgrpc.MetaKeyTraceId, traceId,
	)
	rpcCtx := metadata.NewOutgoingContext(ctx, md)
	return rpcCtx
}
