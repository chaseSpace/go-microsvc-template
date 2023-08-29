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

func (i ClientInterceptor) GRPCCallLog(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	elapsed := time.Now().Sub(start)
	if err != nil {
		errmsg := err.Error()
		if e := xerr.FromErrStr(errmsg); e != nil {
			errmsg = e.FlatMsg()
		}
		xlog.Debug("grpcClient call err log", zap.String("method", method), zap.String("dur", elapsed.String()),
			zap.String("err", errmsg),
			zap.Any("req", req), zap.Any("rsp", reply))
	} else {
		xlog.Debug("grpcClient call log", zap.String("method", method), zap.String("dur", elapsed.String()),
			zap.Any("req", req), zap.Any("rsp", reply))
	}
	return err
}

func (i ClientInterceptor) ExtractGRPCErr(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		e, ok := status.FromError(err)
		if ok {
			err = xerr.ToXErr(errors.New(e.Message()))
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
