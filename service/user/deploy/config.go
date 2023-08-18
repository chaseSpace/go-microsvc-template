package deploy

import (
	"microsvc/deploy"
)

// SvcConfig 每个服务特有的配置结构
type SvcConfig struct {
	Svc      string `mapstructure:"svc"`
	LogLevel string `mapstructure:"log_level"`
}

func (s SvcConfig) GetLogLevel() string {
	return s.LogLevel
}

var _ deploy.SvcConfImpl = new(SvcConfig)

// UserConf 变量命名建议使用服务名作为前缀，避免main文件引用到其他svc的配置变量
var UserConf *SvcConfig
