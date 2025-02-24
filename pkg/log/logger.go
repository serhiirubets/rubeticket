package log

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogrusLogger struct {
	logger *logrus.Logger
}

func (l *LogrusLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *LogrusLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *LogrusLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *LogrusLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *LogrusLogger) WithFields(fields WithFields) ILogger {
	return &LogrusLogger{logger: l.logger.WithFields(logrus.Fields(fields)).Logger}
}

func NewLogrusLogger(logLevel string) ILogger {

	lumberjackLogger := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     14,
		Compress:   true,
	}

	logrusLogger := logrus.New()
	logrusLogger.SetOutput(lumberjackLogger)
	logrusLogger.SetFormatter(&logrus.JSONFormatter{})

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Fatal("Некорректный уровень логирования:", err)
	}

	logrusLogger.SetLevel(level)

	return &LogrusLogger{logger: logrusLogger}
}
