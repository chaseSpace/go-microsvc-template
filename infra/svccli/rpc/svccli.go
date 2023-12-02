//go:build use_sd

// go build -tags use_sd

package rpc

import (
	"google.golang.org/grpc"
	"microsvc/enums/svc"
	"microsvc/infra/svccli"
	"microsvc/protocol/svc/review"
	"microsvc/protocol/svc/user"
)

// If you use this file, Service client use ServiceDiscovery method
// to get service target address, ServiceDiscovery could be implemented
// by Consul/etcd/ZooKeeper/Nacos etc.

var (
	userCli   = svccli.NewCli(svc.User, func(conn *grpc.ClientConn) interface{} { return user.NewUserIntClient(conn) })
	reviewCli = svccli.NewCli(svc.Review, func(conn *grpc.ClientConn) interface{} { return review.NewReviewIntClient(conn) })
)

func User() user.UserIntClient {
	return userCli.Getter().(user.UserIntClient)
}

func Review() review.ReviewIntClient {
	return reviewCli.Getter().(review.ReviewIntClient)
}
