package xtime

import (
	"fmt"
	"time"
)

const (
	Datetime   = "2006-01-02 15:04:05"
	DatetimeMs = "2006-01-02 15:04:05.000"
	Date       = "2006-01-02"
	DateHour   = "2006-01-02 15"
)

func FormatDur(t time.Duration) string {
	if int(t.Seconds()) > 0 {
		return fmt.Sprintf("%.1fs", t.Seconds())
	}
	if t.Milliseconds() > 0 {
		return fmt.Sprintf("%dms", t.Milliseconds())
	}
	if t.Microseconds() > 0 {
		return fmt.Sprintf("%dµs", t.Microseconds())
	}
	return t.String()
}

func GetDayZero(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
