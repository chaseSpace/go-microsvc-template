package rpcext

import (
	"google.golang.org/grpc"
	"microsvc/enums"
	"microsvc/infra/svccli"
	"microsvc/protocol/svc/admin"
	"microsvc/protocol/svc/user"
)

// Service Discovery is now in Multi addresses method.
// TODO: upgrade to DNS method.

var (
	userCli  = svccli.NewCli(enums.SvcUser, func(conn *grpc.ClientConn) interface{} { return user.NewUserExtClient(conn) })
	adminCli = svccli.NewCli(enums.SvcAdmin, func(conn *grpc.ClientConn) interface{} { return admin.NewAdminExtClient(conn) })
)

func User() user.UserExtClient {
	return userCli.Getter().(user.UserExtClient)
}

func Admin() admin.AdminExtClient {
	return adminCli.Getter().(admin.AdminExtClient)
}
