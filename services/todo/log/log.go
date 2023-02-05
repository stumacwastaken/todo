package log

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func init() {
	logConfig := zap.NewProductionConfig()
	logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logConfig.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	logConfig.EncoderConfig.ConsoleSeparator = " "
	logConfig.Encoding = "console"
	logConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	l, err := logConfig.Build()
	if err != nil {
		panic(err)
	}
	log = l
}

func Default() *zap.Logger {
	return log
}

func SetDefault(l *zap.Logger) {
	log = l
}

//New could take an enum level, but given that we're getting this from the flags anyway, I'm not sold on that.
func New(level string) (*zap.Logger, error) {
	atomicLevel := zap.NewAtomicLevel()
	switch strings.ToUpper(level) {
	case "DEBUG":
		atomicLevel.SetLevel(zap.DebugLevel)
	case "INFO":
		atomicLevel.SetLevel(zap.InfoLevel)
	case "WARN":
		atomicLevel.SetLevel(zap.WarnLevel)
	case "ERROR":
		atomicLevel.SetLevel(zap.ErrorLevel)
	default:
		atomicLevel.SetLevel(zap.InfoLevel)
	}
	return newWithLevel(atomicLevel)
}
func newWithLevel(level zap.AtomicLevel) (*zap.Logger, error) {
	logConfig := zap.NewProductionConfig()
	logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logConfig.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	logConfig.EncoderConfig.ConsoleSeparator = " "
	logConfig.Encoding = "console"
	logConfig.Level = level
	return logConfig.Build()
}
