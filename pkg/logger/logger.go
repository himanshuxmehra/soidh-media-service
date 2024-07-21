package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(level string) *zap.Logger {
	config := zap.NewProductionConfig()

	var l zapcore.Level
	err := l.UnmarshalText([]byte(level))
	if err != nil {
		l = zapcore.InfoLevel
	}
	config.Level.SetLevel(l)

	logger, _ := config.Build()
	return logger
}
