package xgrpc

import (
	"context"
	"google.golang.org/grpc/metadata"
)

const (
	MetaKeyFromGateway = "from-gateway"   // store bool to flag grpc request if is from gateway
	MetaKeyAuth        = "authentication" // store token for authentication
	MetaKeyTraceId     = "trace-id"       // store trace id
)

func TransferMetadataWithinCtx(ctx context.Context, key ...string) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		md2 := metadata.New(nil)
		for _, k := range key {
			if len(md[k]) > 0 {
				md2[k] = md[k]
			}
		}
		return metadata.NewOutgoingContext(ctx, md2)
	}
	return ctx
}

func GetMetaVal(ctx context.Context, key string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if ss := md.Get(key); len(ss) > 0 {
			return ss[0]
		}
	}
	return ""
}
