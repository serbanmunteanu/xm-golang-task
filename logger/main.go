package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func SetupErrorLog() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{
		FieldMap: log.FieldMap{
			log.FieldKeyTime: "datetime",
		},
	})
	log.SetLevel(log.InfoLevel)
}

func SetupAccessLog() *log.Logger {
	accessLog := log.New()

	accessLog.SetOutput(os.Stdout)
	accessLog.SetFormatter(&log.JSONFormatter{
		FieldMap: log.FieldMap{
			log.FieldKeyTime: "datetime",
		},
	})
	accessLog.SetLevel(log.InfoLevel)

	return accessLog
}
