package svc

type Svc string

func (s Svc) Name() string {
	if s == "" {
		return "○"
	}
	return "go-" + string(s)
}

const (
	SvcGateway Svc = "gateway"
	SvcUser    Svc = "user"
	SvcAdmin   Svc = "admin"
)
