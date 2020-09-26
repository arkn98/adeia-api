package log

// Logger is the interface for all the functions of a logger.
type Logger interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Sync() error
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
}
