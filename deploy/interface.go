package deploy

import (
	"microsvc/enums"
)

type SvcListenPortSetter interface {
	SetGRPC(int)
	SetHTTP(int)
}

type RegisterSvc interface {
	RegGRPCBase() (name string, addr string, port int)
	RegGRPCMeta() map[string]string
}

type SvcConfImpl interface {
	GetLogLevel() string
	OverrideLogLevel(string)
	GetSvc() enums.Svc
}

type CommConfig struct {
	Svc      enums.Svc `mapstructure:"svc"`
	LogLevel string    `mapstructure:"log_level"`
}

func (s *CommConfig) GetSvc() enums.Svc {
	return s.Svc
}

func (s *CommConfig) GetLogLevel() string {
	return s.LogLevel
}

func (s *CommConfig) OverrideLogLevel(lv string) {
	s.LogLevel = lv
}
