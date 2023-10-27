package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

func Initialize(level string) *logrus.Logger {
	log := logrus.New()
	log.Out = os.Stdout
	switch level {
	case "debug":
		log.Level = logrus.DebugLevel
	case "info":
		log.Level = logrus.InfoLevel
	default:
		log.Level = logrus.ErrorLevel
	}
	log.SetFormatter(&logrus.JSONFormatter{PrettyPrint: true})
	return log
}
