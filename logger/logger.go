package logger

import (
	"github.com/sirupsen/logrus"
)

type Logger struct {
	serviceName string
	logger      *logrus.Logger
	level       logrus.Level
}

func Initialize(serviceName string, level string) *Logger {
	l := &Logger{}
	l.logger = logrus.New()
	l.logger.SetFormatter(&logrus.JSONFormatter{})

	l.serviceName = serviceName
	if level == "" {
		level = "info"
	}

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.InfoLevel
	}

	l.level = lvl
	l.logger.SetLevel(lvl)

	return l
}

func (l *Logger) GetServiceName() string {
	return l.serviceName
}

func (l *Logger) GetLevel() logrus.Level {
	return l.level
}

func (l *Logger) SetLevel(level logrus.Level) {
	l.level = level
	l.logger.SetLevel(level)
}

func (l *Logger) Close() {
	l.logger = nil
}

func (l *Logger) Info(args ...interface{}) {
	l.logger.WithFields(logrus.Fields{
		"service": l.serviceName,
	}).Info(args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.logger.WithFields(logrus.Fields{
		"service": l.serviceName,
	}).Infof(format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.logger.WithFields(logrus.Fields{
		"service": l.serviceName,
	}).Error(args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logger.WithFields(logrus.Fields{
		"service": l.serviceName,
	}).Errorf(format, args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.logger.WithFields(logrus.Fields{
		"service": l.serviceName,
	}).Debug(args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logger.WithFields(logrus.Fields{
		"service": l.serviceName,
	}).Debugf(format, args...)
}

func (l *Logger) Trace(args ...interface{}) {
	l.logger.WithFields(logrus.Fields{
		"service": l.serviceName,
	}).Trace(args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.logger.WithFields(logrus.Fields{
		"service": l.serviceName,
	}).Warn(args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.logger.WithFields(logrus.Fields{
		"service": l.serviceName,
	}).Warnf(format, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.logger.WithFields(logrus.Fields{
		"service": l.serviceName,
	}).Fatal(args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logger.WithFields(logrus.Fields{
		"service": l.serviceName,
	}).Fatalf(format, args...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.logger.WithFields(logrus.Fields{
		"service": l.serviceName,
	}).Panic(args...)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.logger.WithFields(logrus.Fields{
		"service": l.serviceName,
	}).Panicf(format, args...)
}

func (l *Logger) Tracef(format string, args ...interface{}) {
	l.logger.WithFields(logrus.Fields{
		"service": l.serviceName,
	}).Tracef(format, args...)
}
