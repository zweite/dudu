package log

import (
	"io"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	SetLevel(logrus.Level)

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warnln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})
}

func NewLogger(out io.Writer, level logrus.Level) Logger {
	logger := logrus.New()
	logger.Out = out
	logger.Level = level
	return logger
}

func ParseLogLevel(logLevel string, defaultLevel string, logger Logger) Logger {
	if logLevel != "" {
		if level, err := logrus.ParseLevel(logLevel); err == nil {
			logger.SetLevel(level)
		}
	} else {
		if level, err := logrus.ParseLevel(defaultLevel); err == nil {
			logger.SetLevel(level)
		}
	}
	return logger
}
