package log

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

const LogLevelKey = "LOGGING_LEVEL"

func init() {
	levelVal := os.Getenv(LogLevelKey)
	var level logrus.Level
	var err error
	if levelVal == "" {
		level = logrus.InfoLevel
	} else {
		level, err = logrus.ParseLevel(levelVal)
	}

	if err != nil {
		fmt.Printf("Logger failed to init, panicing. Err: %+v", err)
	}

	logrus.SetLevel(level)
}

func DebugF(format string, args ...interface{}) {
	logrus.Debugf(format, args)
}

func InfoF(format string, args ...interface{}) {
	logrus.Infof(format, args)
}

func PrintF(format string, args ...interface{}) {
	logrus.Printf(format, args)
}

func WarnF(format string, args ...interface{}) {
	logrus.Warnf(format, args)
}

func ErrorF(format string, args ...interface{}) {
	logrus.Errorf(format, args)
}

func TraceF(format string, args ...interface{}) {
	logrus.Tracef(format, args)
}

func FatalF(format string, args ...interface{}) {
	logrus.Fatalf(format, args)
}

func PanicF(format string, args ...interface{}) {
	logrus.Panicf(format, args)
}
