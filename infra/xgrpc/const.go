package xgrpc

const (
	certRootCN   = "x.microsvc"
	certClientCN = "client.microsvc"
	certServerCN = "server.microsvc"
)
const unmarshalReqErrPrefix = "grpc: error unmarshalling request: "

func specialClientAuth(svc string, dnsNames []string) bool {
	for _, domain := range dnsNames {
		if domain == svc+"."+certClientCN {
			return true
		}
	}
	return false
}
