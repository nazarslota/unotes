package utils

import (
	"os"
	"os/signal"
)

// GracefulShutdown returns a channel(<-chan struct{}) that closes when an attempt is made to terminate a program.
func GracefulShutdown() <-chan struct{} {
	shutdown := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)
		signal.Notify(s, os.Interrupt)

		<-s
		close(shutdown)
	}()
	return shutdown
}
