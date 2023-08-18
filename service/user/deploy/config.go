package deploy

import (
	"microsvc/deploy"
	"microsvc/util"
)

type SvcConfig struct {
	Svc      string `mapstructure:"svc"`
	LogLevel string `mapstructure:"log_level"`
}

func (s SvcConfig) GetLogLevel() string {
	return s.LogLevel
}

var _ deploy.SvcConfImpl = new(SvcConfig)

var UserConf *SvcConfig

func MustSetup(initializers ...deploy.Initializer) {
	deploy.Init("user", UserConf, initializers...)

	util.AssertNotNil(UserConf)
}
