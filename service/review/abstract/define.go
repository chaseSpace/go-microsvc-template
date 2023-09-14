package abstract

import "microsvc/protocol/svc/review"

// ReviewMethod 在接入第三方审核服务时实现这个接口
type ReviewMethod interface {
	ReviewText(uid int64, sex int32, content string, channel string) (ReviewResult, error)
	ReviewImage(uid int64, sex int32, uri string, channel string) (ReviewResult, error)
	ReviewAudio(uid int64, sex int32, uri string, channel string) (ReviewResult, error)
	ReviewVideo(uid int64, sex int32, uri string, channel string) (ReviewResult, error)
}

type ReviewResult struct {
	State           review.ResultState
	ReqId           string
	RiskLabel       string
	RiskDescription string
}
