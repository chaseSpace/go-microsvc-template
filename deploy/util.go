package deploy

import (
	"microsvc/enums"
	"os"
)

// 读取系统环境变量，默认dev
// export MICROSVC_ENV=beta/prod
func readEnv() enums.Environment {
	env := enums.Environment(os.Getenv(enums.EnvVariable))
	switch env {
	case "":
		return enums.EnvDev
	case enums.EnvBeta, enums.EnvDev, enums.EnvProd:
		return env
	default:
		panic("no valid env provided!")
	}
}
