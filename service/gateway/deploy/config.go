package deploy

import (
	"microsvc/deploy"
)

// SvcConfig 每个服务特有的配置结构
type SvcConfig struct {
	deploy.CommConfig `mapstructure:"root"`
	HttpPort          int `mapstructure:"http_port"`
}

var _ deploy.SvcConfImpl = new(SvcConfig)

var GatewayConf = new(SvcConfig)
