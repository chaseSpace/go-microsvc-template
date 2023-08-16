package consts

type Environment string

const (
	EnvDev  Environment = "dev"
	EnvBeta Environment = "beta"
	EnvProd Environment = "prod"
)

const (
	EnvVariable = "MICROSVC_ENV"
)
