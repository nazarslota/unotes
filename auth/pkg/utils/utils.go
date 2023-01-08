package utils

import (
	"net"
	"net/url"
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

// BuildMongoURI builds new MongoDB URI with the given parameters.
func BuildMongoURI(host, port, username, password string) (string, error) {
	proto := "mongodb"
	if len(port) == 0 {
		proto += "+srv"
	}

	var userinfo *url.Userinfo
	if len(username) > 0 && len(password) > 0 {
		userinfo = url.UserPassword(username, password)
	} else if len(username) > 0 {
		userinfo = url.User(username)
	}

	if len(port) > 0 {
		host = net.JoinHostPort(host, port)
	}

	uri := url.URL{
		Scheme: proto,
		User:   userinfo,
		Host:   host,
	}
	return uri.String(), nil
}
