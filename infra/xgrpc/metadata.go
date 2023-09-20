package xgrpc

import (
	"context"
	"google.golang.org/grpc/metadata"
)

const (
	MdKeyAuthToken       = "authorization" // store token for authentication
	MdKeyTraceId         = "trace-id"      // store trace id
	MdKeyFromGatewayFlag = "from-gateway"
)

const MdKeyFlagExist = "1"

func transferMetadataWithinCtx(ctx context.Context, method string, key ...string) (context.Context, error) {
	ctx, err := sutil.setupCtx(ctx, method)
	if err != nil {
		return nil, err
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		md2 := metadata.New(nil)
		for _, k := range key {
			if len(md[k]) > 0 {
				md2[k] = md[k]
			}
		}
		return metadata.NewOutgoingContext(ctx, md2), nil
	}
	return ctx, nil
}

// GetOutgoingMdVal should be used in client side
func GetOutgoingMdVal(ctx context.Context, key string) string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		if ss := md.Get(key); len(ss) > 0 {
			return ss[0]
		}
	}
	return ""
}

// GetIncomingMdVal should be used in server side
func GetIncomingMdVal(ctx context.Context, key string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if ss := md.Get(key); len(ss) > 0 {
			return ss[0]
		}
	}
	return ""
}
