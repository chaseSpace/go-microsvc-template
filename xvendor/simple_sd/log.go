package simple_sd

import (
	"log"
	"sync"
)

const (
	LogLevelDebug = iota
	LogLevelInfo
	LogLevelError
)

type Logger struct {
	level int
}

var Sdlogger Logger
var onceTips sync.Once

func SetLogLevel(level int) {
	Sdlogger.level = level
}

const logPrefix = "simple_sd - "

func (m *Logger) Debug(msg string, args ...interface{}) {
	onceTips.Do(func() {
		log.Printf(logPrefix + "[Debug]: you can call SetLogLevel() to adjust log level\n")
	})
	if m.level == LogLevelDebug {
		log.Printf(logPrefix+"[Debug]: "+msg+"\n", args...)
	}
}

func (m *Logger) Info(msg string, args ...interface{}) {
	if m.level >= LogLevelInfo {
		log.Printf(logPrefix+"[Info]: "+msg+"\n", args...)
	}
}

func (m *Logger) Error(msg string, args ...interface{}) {
	if m.level >= LogLevelError {
		log.Printf(logPrefix+"[Error]: "+msg+"\n", args...)
	}
}
