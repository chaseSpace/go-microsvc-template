package deploy

import (
	"fmt"
	"microsvc/consts"
	"microsvc/enums"
	"os"
)

// 读取系统环境变量，默认dev
// export MICRO_SVC_ENV=dev/beta/prod
func readEnv() enums.Environment {
	env := enums.Environment(os.Getenv(consts.EnvVar))
	switch env {
	case "":
		return enums.EnvDev
	case enums.EnvBeta, enums.EnvDev, enums.EnvProd:
		return env
	default:
		panic(fmt.Sprintf("invalid env provided: %s", env))
	}
}

func readLogLevelFromEnvVar() string {
	lv := os.Getenv(consts.EnvVarLogLevel)
	switch lv {
	case "debug", "info", "warn", "error":
		return lv
	}
	return ""
}
