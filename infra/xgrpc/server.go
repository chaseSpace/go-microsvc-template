package xgrpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/k0kubun/pp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"microsvc/bizcomm/auth"
	"microsvc/deploy"
	"microsvc/pkg/xerr"
	"microsvc/pkg/xlog"
	"microsvc/pkg/xtime"
	proto2 "microsvc/proto"
	"microsvc/protocol/svc"
	"microsvc/util"
	"microsvc/util/graceful"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"
)

type grpcHTTPRegister func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error

// server动态使用grpc端口范围
const grpcPortMin = 60000
const grpcPortMax = 60999

const httpPortMin = 61000
const httpPortMax = 61999

type XgRPC struct {
	svr                              *grpc.Server
	extHttpRegister, intHttpRegister grpcHTTPRegister
}

func New(interceptors ...grpc.UnaryServerInterceptor) *XgRPC {
	svr := newGRPCServer(deploy.XConf.Svc.Name(), interceptors...)
	return &XgRPC{
		svr:             svr,
		extHttpRegister: nil,
	}
}

func (x *XgRPC) Apply(regFunc func(s *grpc.Server)) {
	regFunc(x.svr)
}

func (x *XgRPC) SetHTTPExtRegister(register grpcHTTPRegister) {
	x.extHttpRegister = register
}

func (x *XgRPC) SetHTTPIntRegister(register grpcHTTPRegister) {
	x.intHttpRegister = register
}

func (x *XgRPC) Start(portSetter deploy.SvcListenPortSetter) {
	lisFetcher := util.NewTcpListenerFetcher(grpcPortMin, grpcPortMax)
	lis, port, err := lisFetcher.Get()
	if err != nil {
		xlog.Panic("failed to get grpc listener", zap.Error(err))
	}
	portSetter.SetGRPC(port)
	grpcAddr := fmt.Sprintf("localhost:%d", port)

	fmt.Printf("\nCongratulations! ^_^\n")
	_, _ = pp.Printf("Your service [%s] is serving gRPC on %s\n", portSetter.GetSvc(), grpcAddr)

	defer graceful.AddStopFunc(func() { // grpc server should stop before http
		x.svr.GracefulStop()
		xlog.Info("xgrpc: gRPC server shutdown completed")
	})

	graceful.Schedule(func() {
		err = x.svr.Serve(lis)
		if err != nil {
			xlog.Error("xgrpc: failed to serve GRPC", zap.String("grpcAddr", grpcAddr), zap.Error(err))
		}
	})

	// 可能需要为grpc服务添加HTTP代理网关
	// NOTE：如果是gateway架构，则不需要
	if x.extHttpRegister != nil || x.intHttpRegister != nil {
		lisFetcher = util.NewTcpListenerFetcher(httpPortMin, httpPortMax)
		lis, port, err := lisFetcher.Get()
		if err != nil {
			xlog.Panic("failed to get http listener", zap.Error(err))
		}
		portSetter.SetHTTP(port)
		httpAddr := fmt.Sprintf(":%d", port)
		fmt.Printf("serving HTTP on http://localhost%s\n", httpAddr)

		graceful.Schedule(func() {
			time.Sleep(time.Second)
			serveHTTP(grpcAddr, lis, x.extHttpRegister, x.intHttpRegister)
		})
	}
	fmt.Println()
}

func serveHTTP(grpcAddr string, httpListener net.Listener, extHandlerRegister, intHandlerRegister grpcHTTPRegister) {
	conn, err := grpc.Dial(grpcAddr, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		xlog.Panic("xgrpc: grpc.Dial failed", zap.String("grpcAddr", grpcAddr), zap.Error(err))
	}
	defer conn.Close()

	muxOpt := newHTTPMuxOpts()
	mux := runtime.NewServeMux(muxOpt...) // create http gateway router for grpc service

	if extHandlerRegister != nil {
		err = extHandlerRegister(context.TODO(), mux, conn)
		if err != nil {
			xlog.Panic("xgrpc: register ext handler failed", zap.String("grpcAddr", grpcAddr), zap.Error(err))
		}
	}
	if intHandlerRegister != nil {
		err = intHandlerRegister(context.TODO(), mux, conn)
		if err != nil {
			xlog.Panic("xgrpc: register int handler failed", zap.String("grpcAddr", grpcAddr), zap.Error(err))
		}
	}
	svr := http.Server{Handler: mux}
	graceful.AddStopFunc(func() {
		util.RunTaskWithCtxTimeout(time.Second*3, func(ctx context.Context) {
			err = svr.Shutdown(ctx)
			xlog.Info("xgrpc: HTTP server shutdown completed", zap.Error(err))
		})
	})

	err = svr.Serve(httpListener)
	if err != nil && err != http.ErrServerClosed {
		xlog.Panic("xgrpc: failed to serve HTTP", zap.String("grpcAddr", grpcAddr), zap.Error(err))
	}
}

type proxyRespMarshaler struct {
	runtime.JSONPb
}

func (c *proxyRespMarshaler) Marshal(grpcRsp interface{}) (b []byte, err error) {
	lastResp := &svc.HttpCommonRsp{
		Code: xerr.ErrNil.Code,
		Msg:  xerr.ErrNil.Msg,
		Data: nil,
	}
	defer func() {
		b, err = c.JSONPb.Marshal(lastResp)
	}()
	if grpcRsp == nil {
		lastResp.Code = xerr.ErrInternal.Code
		lastResp.Msg = "http-proxy: no error, but grpc response is empty"
		return
	}
	data, err := anypb.New(grpcRsp.(proto.Message))
	if err != nil {
		lastResp.Code = xerr.ErrInternal.Code
		lastResp.Msg = fmt.Sprintf("http-proxy: call anypb.New() failed: %v, rsp:%+v", err, grpcRsp)
		return
	}
	lastResp.Data = data
	return
}

func gatewayMarshaler() *proxyRespMarshaler {
	jpb := runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			AllowPartial:    true,
			UseProtoNames:   true,
			UseEnumNumbers:  true,
			EmitUnpopulated: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			AllowPartial:   true,
			DiscardUnknown: true,
		},
	}
	return &proxyRespMarshaler{JSONPb: jpb}
}

func newHTTPMuxOpts() []runtime.ServeMuxOption {
	marshaler := gatewayMarshaler()
	return []runtime.ServeMuxOption{
		runtime.WithMarshalerOption(marshaler.ContentType(nil), marshaler),
		runtime.WithIncomingHeaderMatcher(func(s string) (string, bool) {
			var header = map[string]bool{
				"x-token": true,
			}
			s = strings.ToLower(s)
			return s, header[s]
		}),
		runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, writer http.ResponseWriter, request *http.Request, err error) {
			rsp := &svc.HttpCommonRsp{
				Code: xerr.ErrInternal.Code,
				Msg:  err.Error(),
			}
			s, ok := status.FromError(err)
			if ok {
				if e, ok := xerr.FromErrStr(s.Message()); ok {
					rsp.Code = e.Code
					rsp.Msg = e.Msg
				} else {
					rsp.Msg = s.Message()
				}
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write(util.ToJson(rsp))
		}),
	}
}

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

// -------- 下面是grpc中间件 -----------

type IsExtApiReq interface {
	GetBase() *svc.BaseExtReq
}

type IsAdminApiReq interface {
	GetBase() *svc.AdminBaseReq
}

func RecoverGRPCRequest(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = xerr.ErrInternal.NewMsg(fmt.Sprintf("panic recovered: %v", r))
			xlog.DPanic("RecoverGRPCRequest", zap.String("method", info.FullMethod), zap.Any("err", r),
				zap.String("trace-id", GetIncomingMdVal(ctx, MdKeyTraceId)))
			fmt.Printf("PANIC %v\n%s", r, string(debug.Stack()))
		}
	}()
	return handler(ctx, req)
}

func ToCommonResponse(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	if err == nil {
		_, ok := req.(IsExtApiReq)
		_, ok2 := req.(IsAdminApiReq)
		if ok || ok2 {
			return proto2.RespondOK(resp), nil
		}
	}
	return resp, err
}

func LogGRPCRequest(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()
	resp, err = handler(ctx, req)
	elapsed := xtime.FormatDur(time.Since(start))

	zapFields := []zap.Field{
		zap.String("method", info.FullMethod), zap.String("dur", elapsed),
		zap.Any("req", req), zap.String("trace-id", GetIncomingMdVal(ctx, MdKeyTraceId)),
	}
	if err != nil {
		errmsg := err.Error()
		if e, ok := xerr.FromErr(err); ok {
			errmsg = e.FlatMsg()
		}
		zapFields = append(zapFields, zap.String("err", errmsg))
		xlog.Error("grpc reply_err log", zapFields...)
	} else {
		zapFields = append(zapFields, zap.Any("resp", resp))
		xlog.Debug("grpc reply_ok log", zapFields...)
	}
	return
}

func TraceGRPC(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// We just need transfer some necessary metadata to next rpc call
	// see: https://golang2.eddycjy.com/posts/ch3/09-grpc-metadata-creds/
	ctx = TransferMetadataWithinCtx(ctx, MdKeyTraceId)
	return handler(ctx, req)
}

func StandardizationGRPCErr(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	if err != nil {
		e, ok := status.FromError(err)
		if ok {
			return nil, xerr.ToXErr(errors.New(e.Message()))
		}
	}
	return
}

type SvcClaims struct {
	Authenticated auth.SvcUser
	jwt.RegisteredClaims
}

type AdminClaims struct {
	Authenticated auth.AdminUser
	jwt.RegisteredClaims
}

func Authentication(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	tokenStr := GetIncomingMdVal(ctx, MdKeyAuth)
	if strings.TrimSpace(tokenStr) == "" || strings.TrimLeft(tokenStr, "Bearer ") == "" {
		return nil, xerr.ErrUnauthorized.AppendMsg("empty token")
	}
	claims := SvcClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(deploy.XConf.SvcTokenSignKey), nil
	})
	if err != nil {
		return nil, xerr.ErrUnauthorized.AppendMsg(err.Error())
	}
	if !token.Valid || !claims.Authenticated.IsValid() {
		return nil, xerr.ErrUnauthorized
	}
	ctx = context.WithValue(ctx, auth.CtxAuthenticated{}, claims.Authenticated)
	resp, err = handler(ctx, req)
	return
}
