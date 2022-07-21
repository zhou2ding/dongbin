package l

import (
	"blog/pkg/v"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
)

const (
	defaultMaxSize    = 128 // MB
	defaultMaxAge     = 60  // day
	defaultMaxBackups = 30  // 个
)

func newLogger(appName string, opt ...zap.Option) (*zap.Logger, error) {
	// log level
	level := new(zapcore.Level)
	if err := level.UnmarshalText([]byte(v.GetViper().GetString("log_level"))); err != nil {
		return nil, err
	}
	// log path
	path := v.GetViper().GetString("log.path")
	// max log size
	maxSize := v.GetViper().GetInt("log.max_size")
	if maxSize == 0 {
		maxSize = defaultMaxSize
	}
	// max log age
	maxAge := v.GetViper().GetInt("log.max_age")
	if maxAge == 0 {
		maxAge = defaultMaxAge
	}
	// max backups
	maxBackups := v.GetViper().GetInt("log.max_backups")
	if maxBackups == 0 {
		maxBackups = defaultMaxBackups
	}

	// writes to file or not
	writeSyncer := make([]zapcore.WriteSyncer, 0)
	if len(path) > 0 {
		writeSyncer = append(writeSyncer, zapcore.AddSync(&lumberjack.Logger{
			Filename:   strings.Trim(path, "/") + "/" + appName + ".log",
			MaxSize:    maxSize,
			MaxAge:     maxAge,     // max save days
			MaxBackups: maxBackups, // max counts of backup file
			Compress:   true,
		}))
	}
	// writes to console or not
	if v.GetViper().GetBool("log.console_enable") {
		writeSyncer = append(writeSyncer, zapcore.AddSync(os.Stdout))
	}

	// builds a development Logger that writes DebugLevel and above logs to standard error in a human-friendly format
	if _, err := zap.NewDevelopment(); err != nil {
		return nil, err
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()), // encoder config
		zapcore.NewMultiWriteSyncer(writeSyncer...),                  // destination of writing log
		zap.NewAtomicLevelAt(*level),                                 // log level
	)
	caller := zap.AddCaller() // write filename, line number, and function name in log
	dev := zap.Development()  // development mode, write panic detail in log

	return zap.New(core, caller, dev), nil
}
