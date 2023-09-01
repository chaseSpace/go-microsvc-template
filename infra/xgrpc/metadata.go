package xgrpc

const (
	MetaKeyFromGateway = "meta-from-gateway"   // store bool to flag grpc request if is from gateway
	MetaKeyAuth        = "meta-authentication" // store token for authentication
	MetaKeyTraceId     = "meta-trace-id"       // store trace id
)
