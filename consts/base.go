package consts

const (
	EnvVariable = "MICROSVC_ENV"
)

type Svc string

func (s Svc) Name() string {
	return "svc-" + string(s)
}

const (
	SvcUser  Svc = "user"
	SvcAdmin Svc = "admin"
)
