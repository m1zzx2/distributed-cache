package log

import "distributed-cache/log/impl"

var (
	// std is the name of the standard logger in stdlib `log`
	std = New()
)

func New() Logger {
	return &impl.StdLogger{}
}

func Infof(format string, args ...interface{}) {
	std.Infof(format, args...)
}

func Errorf(format string, args ...interface{}) {
	std.Errorf(format, args...)
}
