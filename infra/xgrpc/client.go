package xgrpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/k0kubun/pp"
	"github.com/samber/lo"
	"github.com/sony/gobreaker"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"microsvc/deploy"
	"microsvc/pkg/xerr"
	"microsvc/pkg/xlog"
	"microsvc/pkg/xtime"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const invalidAddress = "invalidAddress"
const maxRpcRetry = 2

var circuitBreaker = circuitBreakerT{cmap: make(map[string]*gobreaker.CircuitBreaker)}

type circuitBreakerT struct {
	mu   sync.RWMutex
	cmap map[string]*gobreaker.CircuitBreaker
}

func (c *circuitBreakerT) Get(svc string) *gobreaker.CircuitBreaker {
	c.mu.RLock()
	cb := c.cmap[svc]
	c.mu.RUnlock()
	if cb != nil {
		return cb
	}
	c.mu.Lock()
	c.mu.Unlock()
	cb = cutil.newCircuitBreaker(svc)
	c.cmap[svc] = cb
	return cb
}

func NewInvalidGRPCConn(svc string) *grpc.ClientConn {
	cc, err := grpc.Dial(invalidAddress, grpc.WithInsecure(), withClientInterceptorOpt(svc))
	if err != nil {
		panic(err)
	}
	return cc
}

func NewGRPCClient(target, svc string) (cc *grpc.ClientConn, err error) {
	certDir := filepath.Join(deploy.XConf.GetConfDir(), "cert")

	certPath := filepath.Join(certDir, "client-cert.pem")
	keyPath := filepath.Join(certDir, "client-key.pem")

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		panic(err)
	}

	// 加载根证书池，用于验证服务器证书
	rootCA, err := os.ReadFile(filepath.Join(certDir, "ca-cert.pem"))
	if err != nil {
		panic(err)
	}
	rootCAPool := x509.NewCertPool()
	ok := rootCAPool.AppendCertsFromPEM(rootCA)
	if !ok {
		panic("NewGRPCClient: rootCAPool.AppendCertsFromPEM failed")
	}
	// 创建Client TLS 配置
	// 这里使用根证书对server进行验证

	/* 大致流程：
	1. Client 通过请求得到 Server 端的证书
	2. 使用 CA 认证的根证书对 Server 端的证书进行可靠性、有效性等校验
	3. 校验 ServerName 是否匹配
	*/
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      rootCAPool,

		// 在自定义验证逻辑里面，添加证书过期时告警的逻辑，而不是返回error
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			fmt.Printf("\n")
			defer fmt.Printf("\n")

			// 验证证书链中的每个证书（一般是 服务端证书、根证书的顺序）
			for _, chain := range verifiedChains {
				for _, cert := range chain {
					switch cert.Subject.CommonName {
					case certServerCN:
						//pp.Printf("验证Server证书信息: CN:%s before:%s  after:%s \n",
						//	cert.Subject.CommonName, cert.NotBefore, cert.NotAfter)
					case certRootCN:
						//pp.Printf("验证根证书信息: CN:%s before:%s  after:%s \n",
						//	cert.Subject.CommonName, cert.NotBefore, cert.NotAfter)
					default:
						return fmt.Errorf("grpc: handshake faield, server certificate has invalid CN(%s)", cert.Subject.CommonName)
					}
					// 获取证书的有效期
					now := time.Now()
					if now.Before(cert.NotBefore) {
						return fmt.Errorf("grpc: handshake faield, server certificate is invalid before %s", cert.NotBefore)
					}

					if now.After(cert.NotAfter) {
						// 这一步可以不做强验证，因为一旦证书过期（忘记及时更新），这里返回err会导致服务间通信失败
						// 这里可以加上告警
						//return fmt.Errorf("server certificate is expired at %s", cert.NotAfter)

						pp.Printf("server certificate is expired at %s", cert.NotAfter)
					}
				}
			}
			return nil
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// 创建gRPC连接
	cc, err = grpc.DialContext(ctx, target,
		withClientInterceptorOpt(svc),
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
	)
	return
}

type ClientInterceptor struct {
	svc string
}

func newClientInterceptor(svc string) ClientInterceptor {
	return ClientInterceptor{svc: svc}
}

func withClientInterceptorOpt(svc string) grpc.DialOption {
	inter := newClientInterceptor(svc)
	return grpc.WithChainUnaryInterceptor(
		inter.GRPCCallLog,
		inter.ExtractGRPCErr,
		inter.CircuitBreaker,
		inter.Retry,
		inter.WithFailedClient,
	) // execute in reverse order
}

func (i ClientInterceptor) GRPCCallLog(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()

	var err error

	defer func() {
		elapsed := xtime.FormatDur(time.Since(start))
		zapFields := []zap.Field{
			zap.String("method", method), zap.String("dur", elapsed), zap.String("trace-id", GetOutgoingMdVal(ctx, MdKeyTraceId)),
			zap.Any("req", req), zap.Any("rsp", reply),
		}

		req, reply = cutil.beautifyReqAndResInClient(ctx, req, reply)
		if err != nil {
			errmsg := err.Error()
			if e, ok := xerr.FromErrStr(errmsg); ok {
				errmsg = e.FlatMsg()
			}
			xlog.Error("grpc call_err", append(zapFields, zap.String("err", errmsg))...)
		} else {
			xlog.Info("grpc call_ok", zapFields...)
		}
	}()

	err = invoker(ctx, method, req, reply, cc, opts...)
	return err
}

func (i ClientInterceptor) ExtractGRPCErr(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		e, ok := status.FromError(err)
		if ok {
			if e.Message() == context.DeadlineExceeded.Error() {
				return xerr.ErrRPCTimeout
			}
			if strings.HasPrefix(e.Message(), grpcUnmarshalReqErrPrefix) {
				return xerr.ErrBadRequest.AppendMsg(method).AppendMsg(e.Message()[len(grpcUnmarshalReqErrPrefix):])
			}
			err = xerr.ToXErr(errors.New(e.Message()))
		} else {
			err = xerr.ToXErr(err)
		}
	}
	return err
}

// CircuitBreaker executes after retry interceptor
func (i ClientInterceptor) CircuitBreaker(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
	_, cerr := circuitBreaker.Get(i.svc).Execute(func() (interface{}, error) {
		err = invoker(ctx, method, req, reply, cc, opts...)
		if err == nil {
			return nil, nil
		}
		if cutil.breakerTakeError(err) {
			return nil, err
		}
		// ignore other errors
		return nil, nil
	})
	if cerr != nil {
		return cerr
	}
	return
}

func (i ClientInterceptor) Retry(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
	for i := 0; i < maxRpcRetry; i++ {
		err = invoker(ctx, method, req, reply, cc, opts...)
		if err == nil {
			return nil
		}
		// retry it on timeout
		if s, ok := status.FromError(err); ok && lo.Contains([]codes.Code{codes.Unavailable, codes.DeadlineExceeded}, s.Code()) {
			continue
		}
		return err
	}
	return
}

func (i ClientInterceptor) WithFailedClient(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	if cc.Target() == invalidAddress {
		return xerr.ErrServiceUnavailable.AppendMsg("%s", i.svc)
	}
	return invoker(ctx, method, req, reply, cc, opts...)
}
