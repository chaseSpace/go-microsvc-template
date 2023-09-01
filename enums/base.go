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
