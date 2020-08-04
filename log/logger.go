package log

import "github.com/sirupsen/logrus"

// simply logger
type SimpleLogger interface {
	Info(...interface{})
	Infof(string, ...interface{})
	Warn(...interface{})
	Warnf(string, ...interface{})
	Debug(...interface{})
	Debugf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
}

func NewSimpleLogger() SimpleLogger {
	l := logrus.New()
	return l
}
