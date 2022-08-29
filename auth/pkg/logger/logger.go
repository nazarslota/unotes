package logger

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

var instance *log.Logger

func init() {
	instance = log.New()
	instance.SetOutput(os.Stdout)
	instance.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.RFC1123,
	})
}

// L returns logrus instance.
func L() *log.Logger {
	return instance
}
