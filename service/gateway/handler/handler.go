package handler

import (
	"context"
	"fmt"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"microsvc/enums"
	"microsvc/infra/svccli"
	"microsvc/infra/xgrpc/proto"
	"microsvc/pkg/xerr"
	"microsvc/pkg/xlog"
	"microsvc/protocol/svc"
	"microsvc/util"
	"strings"
	"time"
)

type GatewayCtrl struct {
}

const applicationJson = "application/json"

func (GatewayCtrl) Handler(fctx *fasthttp.RequestCtx) {
	errtyp := xerr.ErrOK
	var (
		req         *svc.ForwardReq
		res         = &svc.ForwardRes{}
		fromGateWay = true
	)
	path := string(fctx.Path())

	defer func() {
		fctx.SetContentType(applicationJson)
		fctx.SetStatusCode(200)

		if errtyp.IsOK() {
			xlog.Info("Handler_OK", zap.String("path", path))
			fctx.SetBody(res.Body) // transparent forwarding body
		} else {
			xlog.Info("Handler_FAIL", zap.String("path", path), zap.Error(errtyp))
			httpRes := &svc.GatewayHttpRsp{Code: errtyp.Code, Msg: errtyp.Msg, FromGateway: fromGateWay}
			fctx.SetBody(util.ToJson(httpRes))
		}
	}()

	route := parseRoute(strings.TrimLeft(path, "/"))
	if route == nil {
		errtyp = xerr.ErrApiNotFound
		return
	}

	conn := svccli.GetConn(route.Svc)
	if conn == nil {
		errtyp = xerr.ErrNoRPCClient.AppendMsg(route.Svc.Name())
		return
	}

	// below is grpc calling, we set fromGateWay to false whether the call returns an error or not
	fromGateWay = false
	req = &svc.ForwardReq{Method: route.ForwardMethod(), Body: fctx.PostBody()}

	util.RunTaskWithCtxTimeout(time.Second*5, func(ctx context.Context) {
		md := metadata.Pairs("content-type", "json")
		err := conn.Invoke(newRpcCtx(ctx), path, req, res, grpc.Header(&md))
		if err != nil {
			errtyp = err.(xerr.XErr) // err is converted to XErr in grpc client interceptor
		}
	})
}

type SvcApiRoute struct {
	Svc    enums.Svc
	Prefix string
	Method string
}

func (r SvcApiRoute) ForwardMethod() string {
	return fmt.Sprintf("%s/Forward", r.Prefix)
}

// e.g. path is "svc.user/v1/GetUser"
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

func newRpcCtx(ctx context.Context) context.Context {
	rpcCtx := metadata.NewOutgoingContext(ctx, metadata.Pairs("content-type", proto.Json))
	return rpcCtx
}
