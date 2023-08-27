package deploy

import "microsvc/consts"

type RegisterSvc interface {
	RegBase() (name string, addr string, port int)
	RegMeta() map[string]string
}

type SvcConfImpl interface {
	RegisterSvc
	GetLogLevel() string
	GetSvc() consts.Svc
}
