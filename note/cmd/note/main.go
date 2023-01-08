package main

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	httphandler "github.com/udholdenhed/unotes/note/internal/handler/rest"
)

func main() {
	InitLogger()

	httpServer := httphandler.NewHandler(
		httphandler.WithAddress(":8082"),
	).S()

	if err := httpServer.Serve(); err != nil {
		log.WithError(err).Error("HTTP server error.")
	}
}

func InitLogger() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC1123,
	})
	log.SetOutput(os.Stdout)
}
