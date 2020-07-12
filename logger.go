package skeleton

import (
	stdlog "log"
	"os"
)

// standard logger
var log Logger = stdlog.New(os.Stdout, "", 0)

// Logger interface
type Logger interface {
	Println(v ...interface{})
	Printf(format string, a ...interface{})
}

// SetLogger
func SetLogger(l Logger) {
	if l == nil {
		log.Println("[ERR] Logger is nil")
	}
	log = l
}
