package gormlog

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"time"
)

type Logger struct {
	zap           *zap.Logger
	logLevel      logger.LogLevel
	slowThreshold time.Duration
}

func NewGormLogger(zap *zap.Logger) *Logger {
	return &Logger{
		zap:           zap,
		logLevel:      logger.Info,
		slowThreshold: time.Second,
	}
}

func (l *Logger) LogMode(lv logger.LogLevel) logger.Interface {
	l.logLevel = lv
	return l
}

func (l *Logger) Info(ctx context.Context, msg string, values ...interface{}) {
	l.zap.Info("gorm", zap.Any("msg", msg), zap.Any("values", values))
}

func (l *Logger) Warn(ctx context.Context, msg string, values ...interface{}) {
	l.zap.Warn("gorm", zap.Any("msg", msg), zap.Any("values", values))

}

func (l *Logger) Error(ctx context.Context, msg string, values ...interface{}) {
	l.zap.Error("gorm", zap.Any("msg", msg), zap.Any("values", values))

}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {

}
