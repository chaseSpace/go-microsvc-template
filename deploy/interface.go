package deploy

import (
	"microsvc/enums/svc"
)

type SvcListenPortSetter interface {
	GetSvc() string
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
	GetSvc() svc.Svc
}

type CommConfig struct {
	Svc      svc.Svc `mapstructure:"svc"`
	LogLevel string  `mapstructure:"log_level"`
}

func (s *CommConfig) GetSvc() svc.Svc {
	return s.Svc
}

func (s *CommConfig) GetLogLevel() string {
	return s.LogLevel
}

func (s *CommConfig) OverrideLogLevel(lv string) {
	s.LogLevel = lv
}
