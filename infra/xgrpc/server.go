package xgrpc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/k0kubun/pp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"microsvc/deploy"
	"os"
	"path/filepath"
	"time"
)

func newGRPCServer(svc string, interceptors ...grpc.UnaryServerInterceptor) *grpc.Server {
	certDir := filepath.Join(deploy.XConf.GetConfDir(), "cert")

	certPath := filepath.Join(certDir, "server-cert.pem")
	keyPath := filepath.Join(certDir, "server-key.pem")

	// 加载服务器证书和私钥
	serverCert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		panic(err)
	}

	// 加载根证书池，用于验证客户端证书
	rootCA, err := os.ReadFile(filepath.Join(certDir, "ca-cert.pem"))
	if err != nil {
		panic(err)
	}
	rootCAPool := x509.NewCertPool()
	ok := rootCAPool.AppendCertsFromPEM(rootCA)
	if !ok {
		panic("newGRPCServer: rootCAPool.AppendCertsFromPEM failed")
	}

	// 创建服务器 TLS 配置
	// 使用根证书验证client证书
	serverTLSConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    rootCAPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,

		// 在自定义验证逻辑里面，添加证书过期时告警的逻辑，而不是返回error
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			fmt.Printf("\n")
			defer func() {
				fmt.Printf("\n")
			}()

			// 验证证书链中的每个证书（一般是 客户端证书、根证书的顺序）
			for _, chain := range verifiedChains {
				for _, cert := range chain {
					switch cert.Subject.CommonName {
					case certClientCN:
						//pp.Printf("验证通过--Client证书信息: CN:%s before:%s  after:%s \n",
						//	cert.Subject.CommonName, cert.NotBefore, cert.NotAfter)
					case certRootCN:
						//pp.Printf("验证通过--根证书信息: CN:%s before:%s  after:%s \n",
						//	cert.Subject.CommonName, cert.NotBefore, cert.NotAfter)
					default:
						// 授权特定client
						if specialClientAuth(svc, cert.DNSNames) {
							//pp.Printf("验证通过--特定client CN：%s  DNSNames: %+v\n", cert.Subject.CommonName, cert.DNSNames)
						} else {
							return fmt.Errorf("grpc: handshake faield, invalid client certificate with CN(%s)", cert.Subject.CommonName)
						}
					}
					// 获取证书的有效期
					now := time.Now()
					if now.Before(cert.NotBefore) {
						return fmt.Errorf("grpc: handshake faield, client certificate is invalid before %s", cert.NotBefore)
					}
					if now.After(cert.NotAfter) {
						// 这一步可以不做强验证，因为一旦证书过期（忘记及时更新），这里返回err会导致服务间通信失败
						// 这里可以加上告警
						//return fmt.Errorf("client certificate is expired at %s", cert.NotAfter)

						pp.Printf("client certificate is expired at %s", cert.NotAfter)
					}
				}
			}
			return nil
		},
	}

	// 创建 gRPC 服务器
	base := []grpc.UnaryServerInterceptor{RecoverGRPCRequest,
		ToCommonResponse, LogGRPCRequest,
		TraceGRPC, StandardizationGRPCErr,
		Authentication}

	server := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(serverTLSConfig)),
		grpc.ChainUnaryInterceptor(
			append(base, interceptors...)...,
		))
	return server
}
