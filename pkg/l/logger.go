package l

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
	logger, err := newLogger(prefix)
	gLogger = logger
	return err
}

func Logger() DBLogger {
	return gLogger
}
