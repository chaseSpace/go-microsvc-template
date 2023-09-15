package consts

import "time"

type Datetime string

func (t Datetime) Time() (time.Time, error) {
	ti, err := time.ParseInLocation(LongDateLayout, string(t), time.Local)
	if err != nil {
		return time.Time{}, err
	}
	return ti, nil
}

const (
	ShortDateLayout = "2006-01-02"
	LongDateLayout  = "2006-01-02 15:04:05"
)
