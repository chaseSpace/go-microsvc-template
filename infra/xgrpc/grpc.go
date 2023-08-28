package xgrpc

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"log"
	"microsvc/deploy"
	"microsvc/pkg/xerr"
	"microsvc/pkg/xlog"
	"microsvc/protocol/svc"
	"microsvc/util"
	"microsvc/util/graceful"
	"net"
	"net/http"
	"runtime/debug"
	"time"
)

type grpcHTTPRegister func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error

const grpcPortMin = 60000
const grpcPortMax = 60999

const httpPortMin = 61000
const httpPortMax = 61999

type XgRPC struct {
	svr                              *grpc.Server
	extHttpRegister, intHttpRegister grpcHTTPRegister
}

func New(interceptors ...grpc.UnaryServerInterceptor) *XgRPC {
	// 创建一个grpc svr，并配置适当的中间件
	base := []grpc.UnaryServerInterceptor{RecoverGrpcRequest, LogGrpcRequest}
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(
		append(base, interceptors...)...,
	))
	return &XgRPC{
		svr:             server,
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
		log.Fatalf("failed to get grpc listener: %v", err)
	}
	portSetter.SetGRPC(port)
	grpcAddr := fmt.Sprintf(":%d", port)

	fmt.Println("\nCongratulations! ^_^")
	fmt.Printf("gRPC Server is ready on grpc://localhost%v\n", grpcAddr)

	defer graceful.AddStopFunc(func() { // grpc server should stop before http
		x.svr.GracefulStop()
		xlog.Info("xgrpc: gRPC server shutdown completed")
	})

	go func() {
		err = x.svr.Serve(lis)
		if err != nil {
			xlog.Panic("xgrpc: failed to serve GRPC", zap.String("grpcAddr", grpcAddr), zap.Error(err))
		}
	}()

	if x.extHttpRegister != nil || x.intHttpRegister != nil {
		lisFetcher = util.NewTcpListenerFetcher(httpPortMin, httpPortMax)
		lis, port, err := lisFetcher.Get()
		if err != nil {
			log.Fatalf("failed to get http listener: %v", err)
		}
		portSetter.SetHTTP(port)
		httpAddr := fmt.Sprintf(":%d", port)
		fmt.Printf("HTTP Server is ready on http://localhost%s\n", httpAddr)
		go func() {
			time.Sleep(time.Second)
			serveHTTP(grpcAddr, lis, x.extHttpRegister, x.intHttpRegister)
		}()
	}
	fmt.Println()
}

func serveHTTP(grpcAddr string, httpListener net.Listener, extHandlerRegister, intHandlerRegister grpcHTTPRegister) {
	conn, err := grpc.Dial(grpcAddr, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		xlog.Panic("xgrpc: grpc.Dial failed", zap.String("grpcAddr", grpcAddr), zap.Error(err))
	}
	defer conn.Close()

	//marshaler := gatewayMarshaler()
	muxOpt := []runtime.ServeMuxOption{
		//runtime.WithMarshalerOption(marshaler.ContentType(nil), marshaler),
		//runtime.WithIncomingHeaderMatcher(func(s string) (string, bool) {
		//	var header = map[string]bool{
		//		"x-token": true,
		//	}
		//	s = strings.ToLower(s)
		//	return s, header[s]
		//}),
	}
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

func gatewayMarshaler() *runtime.JSONPb {
	return &runtime.JSONPb{
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
}

// -------- 下面是grpc中间件 -----------

func WrapAdminRsp(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	rsp, err := handler(ctx, req)
	if err != nil {
		return nil, err
	}
	lastResp := new(svc.AdminCommonRsp)
	if rsp == nil {
		lastResp.Code = xerr.ErrInternal.Ecode
		lastResp.Msg = "mw: no error, but response is empty"
		return lastResp, nil
	}
	data, err := anypb.New(rsp.(proto.Message))
	if err != nil {
		lastResp.Code = xerr.ErrInternal.Ecode
		lastResp.Msg = fmt.Sprintf("mw: call anypb.New() failed: %v", err)
		return lastResp, nil
	}
	lastResp.Data = data
	return lastResp, nil
}

func RecoverGrpcRequest(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = xerr.ErrInternal.NewMsg(fmt.Sprintf("panic recovered: %v", r))
			xlog.DPanic("RecoverGrpcRequest", zap.String("method", info.FullMethod), zap.Any("err", r))
			fmt.Printf("PANIC %v\n%s", r, string(debug.Stack()))
		}
	}()
	rsp, err := handler(ctx, req)
	return rsp, err
}

func LogGrpcRequest(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	rsp, err := handler(ctx, req)
	return rsp, err
}
