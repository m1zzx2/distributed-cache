package impl

import (
	"fmt"
	stdlog "log"
)

type StdLogger struct {
}

func (l *StdLogger) Infof(format string, args ...interface{}) {
	stdlog.Println("INFO: "+fmt.Sprintf(format, args...))
}

func (l *StdLogger) Errorf(format string, args ...interface{}) {
	stdlog.Println("ERROR: "+fmt.Sprintf(format, args...))
}
