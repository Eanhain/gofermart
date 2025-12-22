package logger

import (
	"go.uber.org/zap"
)

type (
	Logger struct {
		*zap.SugaredLogger
	}
)

var logger Logger

func InitialLogger() *Logger {
	zapLogger := zap.Must(zap.NewProduction())
	logger = Logger{zapLogger.Sugar()}
	return &logger
}
