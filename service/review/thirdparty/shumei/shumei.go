package shumei

import "microsvc/service/review/abstract"

type Shumei struct {
}

var _ abstract.ReviewMethod = new(Shumei)

func (s Shumei) ReviewText(uid int64, sex int32, content string, channel string) (abstract.ReviewResult, error) {
	//TODO implement me
	panic("implement me")
}

func (s Shumei) ReviewImage(uid int64, sex int32, uri string, channel string) (abstract.ReviewResult, error) {
	//TODO implement me
	panic("implement me")
}

func (s Shumei) ReviewAudio(uid int64, sex int32, uri string, channel string) (abstract.ReviewResult, error) {
	//TODO implement me
	panic("implement me")
}

func (s Shumei) ReviewVideo(uid int64, sex int32, uri string, channel string) (abstract.ReviewResult, error) {
	//TODO implement me
	panic("implement me")
}

func ReviewText(uid int64, sex int32, content string, channel string) (abstract.ReviewResult, error) {
	return Shumei{}.ReviewText(uid, sex, content, channel)
}

func ReviewImage(uid int64, sex int32, content string, channel string) (abstract.ReviewResult, error) {
	return Shumei{}.ReviewImage(uid, sex, content, channel)
}

func ReviewAudio(uid int64, sex int32, content string, channel string) (abstract.ReviewResult, error) {
	return Shumei{}.ReviewAudio(uid, sex, content, channel)
}

func ReviewVideo(uid int64, sex int32, content string, channel string) (abstract.ReviewResult, error) {
	return Shumei{}.ReviewVideo(uid, sex, content, channel)
}
