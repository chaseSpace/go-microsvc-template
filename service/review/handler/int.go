package handler

import (
	"context"
	"microsvc/protocol/svc/review"
)

type ReviewIntCtrl struct {
}

var _ review.UserIntServer = new(ReviewIntCtrl)

func (r ReviewIntCtrl) ReviewResource(ctx context.Context, req *review.ReviewResourceReq) (*review.ReviewResourceRes, error) {
	//TODO implement me
	panic("implement me")
}
