package xgrpc

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"log"
	"microsvc/deploy"
	"microsvc/pkg/xerr"
	"microsvc/protocol/svc"
	"net"
	"net/http"
	"time"
)

const httpPort = ":3200"

type grpcHTTPRegister func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error

type XgRPC struct {
	svr          *grpc.Server
	httpRegister grpcHTTPRegister
}

func New(interceptors ...grpc.UnaryServerInterceptor) *XgRPC {
	// 创建一个grpc svr，并配置适当的中间件
	base := []grpc.UnaryServerInterceptor{RecoverGrpcRequest, LogGrpcRequest}
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(
		append(base, interceptors...)...,
	))
	return &XgRPC{
		svr:          server,
		httpRegister: nil,
	}
}

func (x *XgRPC) Apply(regFunc func(s *grpc.Server)) {
	regFunc(x.svr)
}

func (x *XgRPC) SetHTTPRegister(httpRegister grpcHTTPRegister) {
	x.httpRegister = httpRegister
}

func (x *XgRPC) Stop() {
	x.svr.GracefulStop()
}

func (x *XgRPC) Serve() {
	grpcPort := fmt.Sprintf(":%d", deploy.XConf.GRPCPort)
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Println("\nCongratulations! ^_^")
	fmt.Printf("GRPC Server is listening on grpc://localhost:%s\n", grpcPort)

	if x.httpRegister != nil {
		go func() {
			time.Sleep(time.Second * 2)
			serveHTTP(grpcPort, x.httpRegister)
		}()
	}

	err = x.svr.Serve(lis)
	if err != nil {
		log.Fatalf("failed to Serve: %v", err)
	}
}

func serveHTTP(grpcAddr string, registerHTTP grpcHTTPRegister) {
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect grpc: %v", err)
	}
	defer conn.Close()

	mux := runtime.NewServeMux()
	//opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = registerHTTP(context.TODO(), mux, conn)
	if err != nil {
		log.Fatalf("Failed to register http handler client: %v", err)
	}
	err = http.ListenAndServe(httpPort, mux)
	if err != nil {
		log.Fatalf("Failed to serve http: %v", err)
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
		lastResp.Msg = "no error, but response is empty."
		return lastResp, nil
	}
	data, err := anypb.New(rsp.(proto.Message))
	if err != nil {
		lastResp.Code = xerr.ErrInternal.Ecode
		lastResp.Msg = fmt.Sprintf("call anypb.New() failed: %v", err)
		return lastResp, nil
	}
	lastResp.Data = data
	return lastResp, nil
}

func RecoverGrpcRequest(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = xerr.ErrInternal.NewMsg(fmt.Sprintf("panic recovered: %v", err))
		}
	}()
	rsp, err := handler(ctx, req)
	return rsp, err
}

func LogGrpcRequest(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	rsp, err := handler(ctx, req)
	return rsp, err
}
