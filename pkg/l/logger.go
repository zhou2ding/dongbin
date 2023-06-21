package l

import (
	"go.uber.org/zap"
)

var gLogger *zap.Logger

func InitLogger(prefix string) error {
	logger, err := newLogger(prefix)
	gLogger = logger
	return err
}

func Logger() *zap.Logger {
	return gLogger
}
