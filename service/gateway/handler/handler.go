package handler

import (
	"context"
	"fmt"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"microsvc/enums"
	"microsvc/infra/svccli"
	"microsvc/pkg/xerr"
	"microsvc/pkg/xlog"
	"microsvc/protocol/svc"
	"microsvc/util"
	"strings"
)

type GatewayCtrl struct {
}

const applicationJson = "application/json"

func (GatewayCtrl) Handler(ctx *fasthttp.RequestCtx) {
	errtyp := xerr.ErrOK
	var (
		req         *svc.ForwardReq
		res         = &svc.ForwardRes{}
		fromGateWay = true
	)
	path := string(ctx.Path())

	defer func() {
		ctx.SetContentType(applicationJson)
		ctx.SetStatusCode(200)

		if errtyp.IsOK() {
			xlog.Info("Handler_OK", zap.String("path", path))
			ctx.SetBody(res.Body) // transparent forwarding body
		} else {
			xlog.Info("Handler_FAIL", zap.String("path", path), zap.Error(errtyp))
			httpRes := &svc.GatewayHttpRsp{Code: errtyp.Code, Msg: errtyp.Msg, FromGateway: fromGateWay}
			ctx.SetBody(util.ToJson(httpRes))
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
	req = &svc.ForwardReq{Method: route.Method, Body: ctx.PostBody()}
	err := conn.Invoke(context.TODO(), route.ForwardMethod(), req, res)
	if err != nil {
		errtyp = err.(xerr.XErr) // err is converted to XErr in grpc client interceptor
	}
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
