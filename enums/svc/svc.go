package svc

type Svc string

// Name 请不要随意修改这个方法的逻辑，因为微服务客户端证书中的subjectAltName包含了这个Name
// 一旦修改，客户端证书需要重新申请，否则RPC请求将会失败：tls: bad certificate
// （一种不安全的做法是取消RPC的双向证书验证）
// 参阅：generate_cert_for_svc.md
func (s Svc) Name() string {
	if s == "" {
		return "unknown-svc"
	}
	return string(s)
}

const (
	Gateway Svc = "gateway"
	User    Svc = "user"
	Admin   Svc = "admin"
	Review  Svc = "review"
)
