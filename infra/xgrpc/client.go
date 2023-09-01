package xgrpc

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"microsvc/pkg/xerr"
	"microsvc/pkg/xlog"
	"time"
)

const invalidAddress = "invalidAddress"

func NewInvalidGRPCConn(svc string) *grpc.ClientConn {
	cc, err := grpc.Dial(invalidAddress, grpc.WithInsecure(), withClientInterceptorOpt(svc))
	if err != nil {
		panic(err)
	}
	return cc
}

func NewGRPCClient(target, svc string) (cc *grpc.ClientConn, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	cc, err = grpc.DialContext(ctx, target, grpc.WithInsecure(), withClientInterceptorOpt(svc))
	return
}

type ClientInterceptor struct {
	svc string
}

func newClientInterceptor(svc string) ClientInterceptor {
	return ClientInterceptor{svc: svc}
}

func withClientInterceptorOpt(svc string) grpc.DialOption {
	inter := newClientInterceptor(svc)
	return grpc.WithChainUnaryInterceptor(inter.GRPCCallLog, inter.ExtractGRPCErr, inter.WithFailedClient) // 逆序执行
}

type forwardReply interface {
	GetBody() []byte
}

type forwardReq interface {
	GetMethod() string
	GetBody() []byte
}

func (i ClientInterceptor) GRPCCallLog(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	elapsed := time.Now().Sub(start)
	if err != nil {
		errmsg := err.Error()
		if e, ok := xerr.FromErrStr(errmsg); ok {
			errmsg = e.FlatMsg()
		}
		_req, _ := req.(forwardReq)
		if _req != nil {
			method = _req.GetMethod() // for better logging effect
			req = string(_req.GetBody())
		}
		xlog.Error("GRPCCallLog_ERR", zap.String("method", method), zap.String("dur", elapsed.String()),
			zap.String("err", errmsg),
			zap.Any("req", req), zap.Any("rsp", reply))
	} else {
		res, _ := reply.(forwardReply)
		if res != nil {
			reply = string(res.GetBody()) // for better logging effect
		}
		xlog.Info("GRPCCallLog_OK", zap.String("method", method), zap.String("dur", elapsed.String()),
			zap.Any("req", req), zap.Any("rsp", reply))
	}
	return err
}

func (i ClientInterceptor) ExtractGRPCErr(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		e, ok := status.FromError(err)
		if ok {
			if e.Message() == context.DeadlineExceeded.Error() {
				return xerr.ErrGatewayTimeout
			}
			err = xerr.ToXErr(errors.New(e.Message()))
		} else {
			err = xerr.ToXErr(err)
		}
	}
	return err
}

func (i ClientInterceptor) WithFailedClient(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	if cc.Target() == invalidAddress {
		return xerr.ErrNoRPCClient.AppendMsg("%s", i.svc)
	}
	return invoker(ctx, method, req, reply, cc, opts...)
}
