package deploy

import "microsvc/deploy"

type SvcConfig struct {
	Svc      string `mapstructure:"svc"`
	LogLevel string `mapstructure:"log_level"`
}

func (s SvcConfig) GetLogLevel() string {
	return s.LogLevel
}

var _ deploy.SvcConfImpl = new(SvcConfig)

var AdminConf *SvcConfig
