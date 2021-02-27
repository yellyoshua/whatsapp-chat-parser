package logger

import (
	"log"
	"os"
)

func newLogger(prefix string) *log.Logger {
	return log.New(os.Stdout, prefix, log.LstdFlags)
}

// Fatal print info passed and close program
func Fatal(format string, v ...interface{}) {
	logger := newLogger("fatal - ")
	logger.Fatalf(format, v...)
}

// Info only print passed
func Info(format string, v ...interface{}) {
	logger := newLogger("info - ")
	logger.Printf(format, v...)
}

// CheckError check if contain a error print and close program
func CheckError(message string, e error) {
	if e != nil {
		Fatal("%s -> %s", message, e)
	}
	return
}
