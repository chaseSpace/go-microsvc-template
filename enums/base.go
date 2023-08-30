package enums

type Environment string

const (
	EnvDev  Environment = "dev"
	EnvBeta Environment = "beta"
	EnvProd Environment = "prod"
)

func (e Environment) S() string {
	return string(e)
}

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
