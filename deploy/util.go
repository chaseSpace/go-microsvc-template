package deploy

import (
	"microsvc/consts"
	"os"
)

// 读取系统环境变量，默认dev
// export MICROSVC_ENV=beta/prod
func readEnv() consts.Environment {
	env := consts.Environment(os.Getenv(consts.EnvVariable))
	switch env {
	case "":
		return consts.EnvDev
	case consts.EnvBeta, consts.EnvDev, consts.EnvProd:
		return env
	default:
		panic("no valid env provided!")
	}
}
