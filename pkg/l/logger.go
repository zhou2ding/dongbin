package l

import (
	"blog/pkg/v"
	"github.com/pkg/errors"
)

type DBLogger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Notice(v ...interface{})
	Noticef(format string, v ...interface{})
	Warning(v ...interface{})
	Warningf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Critical(v ...interface{})
	Criticalf(format string, v ...interface{})
	Cat(category string) DBLogger
}

var gLogger DBLogger

func InitLogger(prefix string) error {
	logType := v.GetViper().GetString("log.type")
	switch logType {
	case "zap":
		logger, err := newLogger(prefix)
		if err != nil {
			return err
		}
		gLogger = logger
		return nil
	}
	return errors.New("wrong log type, please check configuration")
}

func Logger() DBLogger {
	return gLogger
}
