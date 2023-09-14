package rpc

import (
	"google.golang.org/grpc"
	"microsvc/enums"
	"microsvc/infra/svccli"
	"microsvc/protocol/svc/review"
	"microsvc/protocol/svc/user"
)

// Service Discover is now in Multi addresses method.
// TODO: upgrade to DNS method.

var (
	userCli   = svccli.NewCli(enums.SvcUser, func(conn *grpc.ClientConn) interface{} { return user.NewUserIntClient(conn) })
	reviewCli = svccli.NewCli(enums.SvcReview, func(conn *grpc.ClientConn) interface{} { return review.NewReviewIntClient(conn) })
)

func User() user.UserIntClient {
	return userCli.Getter().(user.UserIntClient)
}

func Review() review.ReviewIntClient {
	return reviewCli.Getter().(review.ReviewIntClient)
}
