package log

type WithFields map[string]interface{}

type ILogger interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})
	WithFields(fields WithFields) ILogger
}
