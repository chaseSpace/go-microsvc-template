package xgrpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/k0kubun/pp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"microsvc/deploy"
	"microsvc/pkg/xerr"
	"microsvc/pkg/xlog"
	"microsvc/pkg/xtime"
	proto2 "microsvc/proto"
	"microsvc/protocol/svc"
	"microsvc/util"
	"microsvc/util/graceful"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

type grpcHTTPRegister func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error

// server动态使用grpc端口范围
const grpcPortMin = 60000
const grpcPortMax = 60999

const httpPortMin = 61000
const httpPortMax = 61999

type XgRPC struct {
	svr                              *grpc.Server
	extHttpRegister, intHttpRegister grpcHTTPRegister
}

func New(interceptors ...grpc.UnaryServerInterceptor) *XgRPC {
	svr := newGRPCServer(deploy.XConf.Svc.Name(), interceptors...)
	return &XgRPC{
		svr:             svr,
		extHttpRegister: nil,
	}
}

func (x *XgRPC) Apply(regFunc func(s *grpc.Server)) {
	regFunc(x.svr)
}

func (x *XgRPC) SetHTTPExtRegister(register grpcHTTPRegister) {
	x.extHttpRegister = register
}

func (x *XgRPC) SetHTTPIntRegister(register grpcHTTPRegister) {
	x.intHttpRegister = register
}

func (x *XgRPC) Start(portSetter deploy.SvcListenPortSetter) {
	lisFetcher := util.NewTcpListenerFetcher(grpcPortMin, grpcPortMax)
	lis, port, err := lisFetcher.Get()
	if err != nil {
		xlog.Panic("failed to get grpc listener", zap.Error(err))
	}
	portSetter.SetGRPC(port)
	grpcAddr := fmt.Sprintf("localhost:%d", port)

	fmt.Printf("\nCongratulations! ^_^\n")
	_, _ = pp.Printf("Your service [%s] is serving gRPC on %s\n", portSetter.GetSvc(), grpcAddr)

	defer graceful.AddStopFunc(func() { // grpc server should stop before http
		x.svr.GracefulStop()
		xlog.Info("xgrpc: gRPC server shutdown completed")
	})

	graceful.Schedule(func() {
		err = x.svr.Serve(lis)
		if err != nil {
			xlog.Error("xgrpc: failed to serve GRPC", zap.String("grpcAddr", grpcAddr), zap.Error(err))
		}
	})

	// 可能需要为grpc服务添加HTTP代理网关
	// NOTE：如果是gateway架构，则不需要
	if x.extHttpRegister != nil || x.intHttpRegister != nil {
		lisFetcher = util.NewTcpListenerFetcher(httpPortMin, httpPortMax)
		lis, port, err := lisFetcher.Get()
		if err != nil {
			xlog.Panic("failed to get http listener", zap.Error(err))
		}
		portSetter.SetHTTP(port)
		httpAddr := fmt.Sprintf(":%d", port)
		fmt.Printf("serving HTTP on http://localhost%s\n", httpAddr)

		graceful.Schedule(func() {
			time.Sleep(time.Second)
			serveHTTP(grpcAddr, lis, x.extHttpRegister, x.intHttpRegister)
		})
	}
	fmt.Println()
}

func serveHTTP(grpcAddr string, httpListener net.Listener, extHandlerRegister, intHandlerRegister grpcHTTPRegister) {
	conn, err := grpc.Dial(grpcAddr, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		xlog.Panic("xgrpc: grpc.Dial failed", zap.String("grpcAddr", grpcAddr), zap.Error(err))
	}
	defer conn.Close()

	muxOpt := newHTTPMuxOpts()
	mux := runtime.NewServeMux(muxOpt...) // create http gateway router for grpc service

	if extHandlerRegister != nil {
		err = extHandlerRegister(context.TODO(), mux, conn)
		if err != nil {
			xlog.Panic("xgrpc: register ext handler failed", zap.String("grpcAddr", grpcAddr), zap.Error(err))
		}
	}
	if intHandlerRegister != nil {
		err = intHandlerRegister(context.TODO(), mux, conn)
		if err != nil {
			xlog.Panic("xgrpc: register int handler failed", zap.String("grpcAddr", grpcAddr), zap.Error(err))
		}
	}
	svr := http.Server{Handler: mux}
	graceful.AddStopFunc(func() {
		util.RunTaskWithCtxTimeout(time.Second*3, func(ctx context.Context) {
			err = svr.Shutdown(ctx)
			xlog.Info("xgrpc: HTTP server shutdown completed", zap.Error(err))
		})
	})

	err = svr.Serve(httpListener)
	if err != nil && err != http.ErrServerClosed {
		xlog.Panic("xgrpc: failed to serve HTTP", zap.String("grpcAddr", grpcAddr), zap.Error(err))
	}
}

type proxyRespMarshaler struct {
	runtime.JSONPb
}

func (c *proxyRespMarshaler) Marshal(grpcRsp interface{}) (b []byte, err error) {
	lastResp := &svc.HttpCommonRsp{
		Code: xerr.ErrNil.Code,
		Msg:  xerr.ErrNil.Msg,
		Data: nil,
	}
	defer func() {
		b, err = c.JSONPb.Marshal(lastResp)
	}()
	if grpcRsp == nil {
		lastResp.Code = xerr.ErrInternal.Code
		lastResp.Msg = "http-proxy: no error, but grpc response is empty"
		return
	}
	data, err := anypb.New(grpcRsp.(proto.Message))
	if err != nil {
		lastResp.Code = xerr.ErrInternal.Code
		lastResp.Msg = fmt.Sprintf("http-proxy: call anypb.New() failed: %v, rsp:%+v", err, grpcRsp)
		return
	}
	lastResp.Data = data
	return
}

func gatewayMarshaler() *proxyRespMarshaler {
	jpb := runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			AllowPartial:    true,
			UseProtoNames:   true,
			UseEnumNumbers:  true,
			EmitUnpopulated: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			AllowPartial:   true,
			DiscardUnknown: true,
		},
	}
	return &proxyRespMarshaler{JSONPb: jpb}
}

func newHTTPMuxOpts() []runtime.ServeMuxOption {
	marshaler := gatewayMarshaler()
	return []runtime.ServeMuxOption{
		runtime.WithMarshalerOption(marshaler.ContentType(nil), marshaler),
		runtime.WithIncomingHeaderMatcher(func(s string) (string, bool) {
			var header = map[string]bool{
				"x-token": true,
			}
			s = strings.ToLower(s)
			return s, header[s]
		}),
		runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, writer http.ResponseWriter, request *http.Request, err error) {
			rsp := &svc.HttpCommonRsp{
				Code: xerr.ErrInternal.Code,
				Msg:  err.Error(),
			}
			s, ok := status.FromError(err)
			if ok {
				if e, ok := xerr.FromErrStr(s.Message()); ok {
					rsp.Code = e.Code
					rsp.Msg = e.Msg
				} else {
					rsp.Msg = s.Message()
				}
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write(util.ToJson(rsp))
		}),
	}
}

// -------- 下面是grpc中间件 -----------

type IsExtApiReq interface {
	GetBase() *svc.BaseExtReq
}

type IsAdminApiReq interface {
	GetBase() *svc.AdminBaseReq
}

func RecoverGRPCRequest(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = xerr.ErrInternal.NewMsg(fmt.Sprintf("panic recovered: %v", r))
			xlog.DPanic("RecoverGRPCRequest", zap.String("method", info.FullMethod), zap.Any("err", r),
				zap.String("trace-id", GetMetaVal(ctx, MetaKeyTraceId)))
			fmt.Printf("PANIC %v\n%s", r, string(debug.Stack()))
		}
	}()
	return handler(ctx, req)
}

func ToCommonResponse(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	if err == nil {
		_, ok := req.(IsExtApiReq)
		_, ok2 := req.(IsAdminApiReq)
		if ok || ok2 {
			return proto2.RespondOK(resp), nil
		}
	}
	return resp, err
}

func LogGRPCRequest(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()
	resp, err = handler(ctx, req)
	elapsed := xtime.FormatDur(time.Since(start))

	zapFields := []zap.Field{
		zap.String("method", info.FullMethod), zap.String("dur", elapsed),
		zap.Any("req", req), zap.String("trace-id", GetMetaVal(ctx, MetaKeyTraceId)),
	}
	if err != nil {
		errmsg := err.Error()
		if e, ok := xerr.FromErr(err); ok {
			errmsg = e.FlatMsg()
		}
		zapFields = append(zapFields, zap.String("err", errmsg))
		xlog.Error("grpc reply_err log", zapFields...)
	} else {
		zapFields = append(zapFields, zap.Any("resp", resp))
		xlog.Debug("grpc reply_ok log", zapFields...)
	}
	return
}

func TraceGRPC(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// We just need transfer some necessary metadata to next rpc call
	// see: https://golang2.eddycjy.com/posts/ch3/09-grpc-metadata-creds/
	ctx = TransferMetadataWithinCtx(ctx, MetaKeyTraceId)
	return handler(ctx, req)
}

func StandardizationGRPCErr(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	if err != nil {
		e, ok := status.FromError(err)
		if ok {
			return nil, xerr.ToXErr(errors.New(e.Message()))
		}
	}
	return
}
