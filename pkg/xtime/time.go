package xtime

import "time"

const (
	Datetime   = "2006-01-02 15:04:05"
	DatetimeMs = "2006-01-02 15:04:05.000"
	Date       = "2006-01-02"
	DateHour   = "2006-01-02 15"
)

func GetDayZero(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
