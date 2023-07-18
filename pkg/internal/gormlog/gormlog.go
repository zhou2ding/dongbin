package gormlog

import (
	"blog/pkg/l"
	"context"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"time"
)

type Logger struct {
	log           l.DBLogger
	logLevel      logger.LogLevel
	slowThreshold time.Duration
}

func NewGormLogger(log l.DBLogger) *Logger {
	return &Logger{
		log:           log,
		logLevel:      logger.Info,
		slowThreshold: time.Second,
	}
}

func (l *Logger) LogMode(lv logger.LogLevel) logger.Interface {
	l.logLevel = lv
	return l
}

func (l *Logger) Info(ctx context.Context, msg string, values ...interface{}) {
	l.log.Info("gorm", zap.Any("msg", msg), zap.Any("values", values))
}

func (l *Logger) Warn(ctx context.Context, msg string, values ...interface{}) {
	l.log.Warning("gorm", zap.Any("msg", msg), zap.Any("values", values))

}

func (l *Logger) Error(ctx context.Context, msg string, values ...interface{}) {
	l.log.Error("gorm", zap.Any("msg", msg), zap.Any("values", values))

}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.logLevel > logger.Silent {
		elapsed := time.Since(begin)
		switch {
		case err != nil && l.logLevel >= logger.Error:
			sql, rows := fc()
			if rows == -1 {
				l.log.Error("gorm", zap.Error(err), zap.Float64("used time", float64(elapsed.Nanoseconds())/1e6), zap.String("SQL", sql))
			} else {
				l.log.Error("gorm", zap.Error(err), zap.Float64("used time", float64(elapsed.Nanoseconds())/1e6), zap.Int64("rows", rows), zap.String("SQL", sql))
			}
		case elapsed > l.slowThreshold && l.slowThreshold != 0 && l.logLevel >= logger.Warn:
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", l.slowThreshold)
			if rows == -1 {
				l.log.Warning("gorm", zap.String("slowLog", slowLog), zap.Float64("used time", float64(elapsed.Nanoseconds())/1e6), zap.String("SQL", sql))
			} else {
				l.log.Warning("gorm", zap.String("slowLog", slowLog), zap.Float64("used time", float64(elapsed.Nanoseconds())/1e6), zap.Int64("rows", rows), zap.String("SQL", sql))
			}
		case l.logLevel == logger.Info:
			sql, rows := fc()
			if rows == -1 {
				l.log.Info("gorm", zap.Float64("used time", float64(elapsed.Nanoseconds())/1e6), zap.String("SQL", sql))
			} else {
				l.log.Info("gorm", zap.Float64("used time", float64(elapsed.Nanoseconds())/1e6), zap.Int64("rows", rows), zap.String("SQL", sql))
			}
		}
	}
}
