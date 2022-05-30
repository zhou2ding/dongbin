package gormlog

import (
	"context"
	"fmt"
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
	if l.logLevel > logger.Silent {
		elapsed := time.Since(begin)
		switch {
		case err != nil && l.logLevel >= logger.Error:
			sql, rows := fc()
			if rows == -1 {
				l.zap.Error("gorm", zap.Error(err), zap.Float64("used time", float64(elapsed.Nanoseconds())/1e6), zap.String("SQL", sql))
			} else {
				l.zap.Error("gorm", zap.Error(err), zap.Float64("used time", float64(elapsed.Nanoseconds())/1e6), zap.Int64("rows", rows), zap.String("SQL", sql))
			}
		case elapsed > l.slowThreshold && l.slowThreshold != 0 && l.logLevel >= logger.Warn:
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", l.slowThreshold)
			if rows == -1 {
				l.zap.Warn("gorm", zap.String("slowLog", slowLog), zap.Float64("used time", float64(elapsed.Nanoseconds())/1e6), zap.String("SQL", sql))
			} else {
				l.zap.Warn("gorm", zap.String("slowLog", slowLog), zap.Float64("used time", float64(elapsed.Nanoseconds())/1e6), zap.Int64("rows", rows), zap.String("SQL", sql))
			}
		case l.logLevel == logger.Info:
			sql, rows := fc()
			if rows == -1 {
				l.zap.Info("gorm", zap.Float64("used time", float64(elapsed.Nanoseconds())/1e6), zap.String("SQL", sql))
			} else {
				l.zap.Info("gorm", zap.Float64("used time", float64(elapsed.Nanoseconds())/1e6), zap.Int64("rows", rows), zap.String("SQL", sql))
			}
		}
	}
}
