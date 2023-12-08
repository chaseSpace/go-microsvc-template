package xgrpc

const (
	certRootCN   = "x.microsvc"
	certClientCN = "client.microsvc"
	certServerCN = "*.default.svc.cluster.local"
)
const grpcUnmarshalReqErrPrefix = "grpc: error unmarshalling request: "

func specialClientAuth(svc string, dnsNames []string) bool {
	for _, domain := range dnsNames {
		if domain == svc+"."+certClientCN {
			return true
		}
	}
	return false
}

type CtxServerSideKey struct{}

type CtxServerSideVal struct {
	IsExtMethod bool
	FromGateway bool
}
