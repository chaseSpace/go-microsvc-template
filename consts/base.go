package consts

const (
	EnvVariable = "MICROSVC_ENV"
)

type CtxKey struct{}

type CtxValue struct {
}

type Svc string

func (s Svc) Name() string {
	return "go-" + string(s)
}

const (
	SvcUser  Svc = "user"
	SvcAdmin Svc = "admin"
)
