package deploy

import "microsvc/consts"

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
	GetSvc() consts.Svc
}
