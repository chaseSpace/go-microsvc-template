package xlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"microsvc/deploy"
	"microsvc/pkg/xtime"
	"os"
)

var xlogger *zap.Logger

func Init(cc *deploy.XConfig) {
	//level := deploy.XConf.GetSvcConf().GetLogLevel()
	var lv = zapcore.DebugLevel
	switch cc.GetSvcConf().GetLogLevel() {
	case "info":
		lv = zapcore.InfoLevel
	case "warning":
		lv = zapcore.WarnLevel
	case "error":
		lv = zapcore.ErrorLevel
	}

	writer := zapcore.AddSync(os.Stdout) // 写stdout，再用容器收集日志
	core := zapcore.NewCore(getEncoder(lv), writer, lv)
	xlogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

func Stop() {
	_ = xlogger.Sync()
}

func getEncoder(level zapcore.Level) zapcore.Encoder {
	var encoderConf zapcore.EncoderConfig
	if level == zapcore.DebugLevel {
		encoderConf = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderConf = zap.NewProductionEncoderConfig()
	}
	encoderConf.EncodeTime = zapcore.TimeEncoderOfLayout(xtime.Datetime)
	return zapcore.NewConsoleEncoder(encoderConf)
}

// --------------------------------

func Debug(msg string, fields ...zapcore.Field) {
	xlogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zapcore.Field) {
	xlogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zapcore.Field) {
	xlogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zapcore.Field) {
	xlogger.Error(msg, fields...)
}

func Panic(msg string, fields ...zapcore.Field) {
	xlogger.Panic(msg, fields...)
}

// 这里不需要DPanic函数，因为Panic够用，且我们的grpc中间件会捕获panic，并封装包含panic信息的Response
