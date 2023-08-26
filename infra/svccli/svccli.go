package svccli

import (
	"google.golang.org/grpc"
	"microsvc/consts"
	"microsvc/protocol/svc/admin"
	"microsvc/protocol/svc/user"
)

// Service Discovery is now in Multi addresses method.
// TODO: upgrade to DNS method.

var (
	userCli  = newIntCli(consts.SvcUser, func(conn *grpc.ClientConn) interface{} { return user.NewUserIntClient(conn) })
	adminCli = newIntCli(consts.SvcAdmin, func(conn *grpc.ClientConn) interface{} { return admin.NewAdminSvcClient(conn) })
)

func User() user.UserIntClient {
	return userCli.Getter().(user.UserIntClient)
}

func Admin() admin.AdminSvcClient {
	return adminCli.Getter().(admin.AdminSvcClient)
}
