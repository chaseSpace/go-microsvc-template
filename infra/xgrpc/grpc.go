package xgrpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"log"
	"microsvc/pkg/xerr"
	"microsvc/protocol/svc"
	"net"
)

func ServeGRPC(svr *grpc.Server, port ...int) {
	_port := 3000
	if len(port) > 0 {
		_port = port[0]
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", _port)) // 监听在端口 50051
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Println("\nCongratulations! ^_^")
	fmt.Printf("GRPC Server is listening on grpc://localhost:%d\n", _port)

	err = svr.Serve(lis)
	if err != nil {
		log.Fatalf("failed to Serve: %v", err)
	}
}

// -------- 下面是grpc中间件 -----------

func WrapAdminRsp(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	rsp, err := handler(ctx, req)
	if err != nil {
		return nil, err
	}
	lastResp := new(svc.AdminCommonRsp)
	if rsp == nil {
		lastResp.Code = xerr.ErrInternal.Code
		lastResp.Msg = "no error, but response is empty."
		return lastResp, nil
	}
	data, err := anypb.New(rsp.(proto.Message))
	if err != nil {
		lastResp.Code = xerr.ErrInternal.Code
		lastResp.Msg = fmt.Sprintf("call anypb.New() failed: %v", err)
		return lastResp, nil
	}
	lastResp.Data = data
	return lastResp, nil
}
