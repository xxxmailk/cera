package log

import "github.com/sirupsen/logrus"

// simply logger
type SimpleLogger interface {
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Debugf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
}

func NewSimpleLogger() SimpleLogger {
	l := logrus.New()
	return l
}
