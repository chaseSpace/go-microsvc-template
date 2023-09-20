package svc

type Svc string

func (s Svc) Name() string {
	if s == "" {
		return "unknown-svc"
	}
	return "go-" + string(s)
}

const (
	Gateway Svc = "gateway"
	User    Svc = "user"
	Admin   Svc = "admin"
	Review  Svc = "review"
)
