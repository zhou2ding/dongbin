package l

import (
	"blog/pkg/v"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
)

type zLogger struct {
	*zap.SugaredLogger
}

const (
	defaultMaxSize    = 128 // MB
	defaultMaxAge     = 60  // day
	defaultMaxBackups = 30  // 个
)

func newLogger(appName string, opt ...zap.Option) (*zLogger, error) {
	lvTransfer := map[string]string{
		"DEBU":    "debug",
		"ERRO":    "error",
		"WARNING": "warn",

		"NOTI":     "info",
		"NOTICE":   "info",
		"CRIT":     "dpanic",
		"CRITICAL": "dpanic",
	}
	lvStr := v.GetViper().GetString("logger.level")
	if s, ok := lvTransfer[lvStr]; ok {
		lvStr = s
	}
	// log lv
	lv := new(zapcore.Level)
	if err := lv.UnmarshalText([]byte(lvStr)); err != nil {
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
			Filename:   strings.TrimRight(path, "/") + "/" + appName + ".log",
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
		zap.NewAtomicLevelAt(*lv),                                    // log lv
	)
	caller := zap.AddCaller() // write filename, line number, and function name in log
	dev := zap.Development()  // development mode, write panic detail in log

	return &zLogger{zap.New(core, caller, dev).Sugar()}, nil
}

func (z *zLogger) Print(v ...interface{}) {
	z.SugaredLogger.Debug(v...)
}
func (z *zLogger) Printf(format string, v ...interface{}) {
	z.SugaredLogger.Debugf(format, v...)
}

func (z *zLogger) Fatal(v ...interface{}) {
	z.SugaredLogger.Fatal(v...)
}
func (z *zLogger) Fatalf(format string, v ...interface{}) {
	z.SugaredLogger.Fatalf(format, v...)
}

func (z *zLogger) Panic(v ...interface{}) {
	z.SugaredLogger.Panic(v...)
}
func (z *zLogger) Panicf(format string, v ...interface{}) {
	z.SugaredLogger.Panicf(format, v...)
}

func (z *zLogger) Info(v ...interface{}) {
	z.SugaredLogger.Info(v...)
}
func (z *zLogger) Infof(format string, v ...interface{}) {
	z.SugaredLogger.Infof(format, v...)
}

func (z *zLogger) Debug(v ...interface{}) {
	z.SugaredLogger.Debug(v...)
}
func (z *zLogger) Debugf(format string, v ...interface{}) {
	z.SugaredLogger.Debugf(format, v...)
}

func (z *zLogger) Notice(v ...interface{}) {
	z.SugaredLogger.Info(v...)
}
func (z *zLogger) Noticef(format string, v ...interface{}) {
	z.SugaredLogger.Infof(format, v...)
}

func (z *zLogger) Warning(v ...interface{}) {
	z.SugaredLogger.Warn(v...)
}
func (z *zLogger) Warningf(format string, v ...interface{}) {
	z.SugaredLogger.Warnf(format, v...)
}

func (z *zLogger) Error(v ...interface{}) {
	z.SugaredLogger.Error(v...)

}
func (z *zLogger) Errorf(format string, v ...interface{}) {
	z.SugaredLogger.Errorf(format, v...)
}

func (z *zLogger) Critical(v ...interface{}) {
	z.SugaredLogger.Panic(v...)
}
func (z *zLogger) Criticalf(format string, v ...interface{}) {
	z.SugaredLogger.Panicf(format, v...)
}

func (z *zLogger) Cat(category string) DBLogger {
	return z
}
