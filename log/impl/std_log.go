package impl

import (
	"fmt"
	stdlog "log"
)

type StdLogger struct {
}

func (l *StdLogger) Infof(format string, args ...interface{}) {
	if len(args) == 0 {
		stdlog.Println("INFO: " + format)
		return
	}
	stdlog.Println("INFO: " + fmt.Sprintf(format, args...))
}

func (l *StdLogger) Errorf(format string, args ...interface{}) {
	if len(args) == 0 {
		stdlog.Println("ERROR: " + format)
		return
	}
	stdlog.Println("ERROR: " + fmt.Sprintf(format, args...))
}
