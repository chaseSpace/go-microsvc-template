// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.4
// source: svc/review/review.int.proto

package review

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ReviewIntClient is the client API for ReviewInt service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ReviewIntClient interface {
	// 接口名称不要过于笼统，这样不便于搜索
	// good：ReviewResource
	// bad: Review
	ReviewResource(ctx context.Context, in *ReviewResourceReq, opts ...grpc.CallOption) (*ReviewResourceRes, error)
}

type reviewIntClient struct {
	cc grpc.ClientConnInterface
}

func NewReviewIntClient(cc grpc.ClientConnInterface) ReviewIntClient {
	return &reviewIntClient{cc}
}

func (c *reviewIntClient) ReviewResource(ctx context.Context, in *ReviewResourceReq, opts ...grpc.CallOption) (*ReviewResourceRes, error) {
	out := new(ReviewResourceRes)
	err := c.cc.Invoke(ctx, "/svc.user.ReviewInt/ReviewResource", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReviewIntServer is the server API for ReviewInt service.
// All implementations should embed UnimplementedReviewIntServer
// for forward compatibility
type ReviewIntServer interface {
	// 接口名称不要过于笼统，这样不便于搜索
	// good：ReviewResource
	// bad: Review
	ReviewResource(context.Context, *ReviewResourceReq) (*ReviewResourceRes, error)
}

// UnimplementedReviewIntServer should be embedded to have forward compatible implementations.
type UnimplementedReviewIntServer struct {
}

func (UnimplementedReviewIntServer) ReviewResource(context.Context, *ReviewResourceReq) (*ReviewResourceRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReviewResource not implemented")
}

// UnsafeReviewIntServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReviewIntServer will
// result in compilation errors.
type UnsafeReviewIntServer interface {
	mustEmbedUnimplementedReviewIntServer()
}

func RegisterReviewIntServer(s grpc.ServiceRegistrar, srv ReviewIntServer) {
	s.RegisterService(&ReviewInt_ServiceDesc, srv)
}

func _ReviewInt_ReviewResource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReviewResourceReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReviewIntServer).ReviewResource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/svc.user.ReviewInt/ReviewResource",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReviewIntServer).ReviewResource(ctx, req.(*ReviewResourceReq))
	}
	return interceptor(ctx, in, info, handler)
}

// ReviewInt_ServiceDesc is the grpc.ServiceDesc for ReviewInt service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ReviewInt_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "svc.user.ReviewInt",
	HandlerType: (*ReviewIntServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ReviewResource",
			Handler:    _ReviewInt_ReviewResource_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "svc/review/review.int.proto",
}
