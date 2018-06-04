package log

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
)

func init() {
	logrus.SetFormatter(&prefixed.TextFormatter{
		FullTimestamp:    true,
		ForceColors:      true,
		QuoteEmptyFields: true,
		ForceFormatting:  true,
		TimestampFormat:  "15:04:05.000000 MST",
	})
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(os.Stderr)
}

var (
	Debug = logrus.Debug
	Info  = logrus.Info
	Warn  = logrus.Warn
	Error = logrus.Error
	Fatal = logrus.Fatal

	Debugf = logrus.Debugf
	Infof  = logrus.Infof
	Warnf  = logrus.Warnf
	Errorf = logrus.Errorf
	Fatalf = logrus.Fatalf
)

func Panic(args ...interface{}) {
	logrus.Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	logrus.Panicf(format, args...)
}
