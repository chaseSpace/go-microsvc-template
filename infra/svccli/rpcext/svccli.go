//go:build !k8s

package rpcext

import (
	"google.golang.org/grpc"
	"microsvc/enums/svc"
	"microsvc/infra/svccli"
	"microsvc/protocol/svc/admin"
	"microsvc/protocol/svc/user"
)

// If you use this file, Service client use ServiceDiscovery method
// to get service target address, ServiceDiscovery could be implemented
// by Consul/etcd/ZooKeeper/Nacos etc.

var (
	userCli  = svccli.NewCli(svc.User, func(conn *grpc.ClientConn) interface{} { return user.NewUserExtClient(conn) })
	adminCli = svccli.NewCli(svc.Admin, func(conn *grpc.ClientConn) interface{} { return admin.NewAdminExtClient(conn) })
)

func User() user.UserExtClient {
	return userCli.Getter().(user.UserExtClient)
}

func Admin() admin.AdminExtClient {
	return adminCli.Getter().(admin.AdminExtClient)
}
