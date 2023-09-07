package mdns

import "log"

const (
	LogLevelDebug = iota
	LogLevelInfo
)

type mlog struct {
	level int
}

func (m *mlog) Debug(msg string, args ...interface{}) {
	if m.level == LogLevelDebug {
		log.Printf(msg+"\n", args...)
	}
}

func (m *mlog) Info(msg string, args ...interface{}) {
	if m.level >= LogLevelInfo {
		log.Printf(msg+"\n", args...)
	}
}
