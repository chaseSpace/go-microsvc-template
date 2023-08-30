package api

import "microsvc/protocol/svc/user"

type Operation struct {
	Req func() interface{}
}

// 手工注册API
var userApiRegistration = map[string]*Operation{
	"GetUser": &Operation{Req: func() interface{} { return new(user.GetUserReq) }},
}

func LoadUserApi(api string) *Operation {
	return userApiRegistration[api]
}
