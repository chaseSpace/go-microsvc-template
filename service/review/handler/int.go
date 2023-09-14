package handler

import (
	"context"
	"microsvc/pkg/xerr"
	"microsvc/protocol/svc/review"
	"microsvc/service/review/abstract"
	"microsvc/service/review/thirdparty/shumei"
)

type ReviewIntCtrl struct {
}

var _ review.ReviewIntServer = new(ReviewIntCtrl)

func (r ReviewIntCtrl) ReviewResource(ctx context.Context, req *review.ReviewResourceReq) (*review.ReviewResourceRes, error) {
	res := new(review.ReviewResourceRes)
	var rr abstract.ReviewResult
	var err error
	switch req.Type {
	case review.ReviewType_Text:
		rr, err = shumei.ReviewText(req.Uid, req.Sex, req.Content, req.Channel)
	case review.ReviewType_Image:
		rr, err = shumei.ReviewImage(req.Uid, req.Sex, req.Content, req.Channel)
	case review.ReviewType_Audio:
		rr, err = shumei.ReviewAudio(req.Uid, req.Sex, req.Content, req.Channel)
	case review.ReviewType_Video:
		rr, err = shumei.ReviewVideo(req.Uid, req.Sex, req.Content, req.Channel)
	default:
		return nil, xerr.ErrParams.AppendMsg("unsupported review type:%d", req.Type)
	}
	if err != nil {
		return nil, xerr.ErrThirdParty.AppendMsg(err.Error())
	}
	res.State = rr.State
	return res, nil
}
