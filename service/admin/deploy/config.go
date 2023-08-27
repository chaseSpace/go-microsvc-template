package deploy

import (
	"microsvc/consts"
	"microsvc/deploy"
)

type SvcConfig struct {
	Svc      consts.Svc `mapstructure:"svc"`
	LogLevel string     `mapstructure:"log_level"`
}

func (s SvcConfig) GetSvc() consts.Svc {
	return s.Svc
}

func (s SvcConfig) GetLogLevel() string {
	return s.LogLevel
}

var _ deploy.SvcConfImpl = new(SvcConfig)

var AdminConf = new(SvcConfig)
