package simple_sd

import "log"

const (
	LogLevelDebug = iota
	LogLevelInfo
	LogLevelError
)

type Logger struct {
	level int
}

var Sdlogger Logger

func SetLogLevel(level int) {
	Sdlogger.level = level
}

func (m *Logger) Debug(msg string, args ...interface{}) {
	if m.level == LogLevelDebug {
		log.Printf("[Debug]: "+msg+"\n", args...)
	}
}

func (m *Logger) Info(msg string, args ...interface{}) {
	if m.level >= LogLevelInfo {
		log.Printf("[Info]: "+msg+"\n", args...)
	}
}

func (m *Logger) Error(msg string, args ...interface{}) {
	if m.level >= LogLevelError {
		log.Printf("[Error]: "+msg+"\n", args...)
	}
}
