package utils

import (
	"errors"
	"net"
	"net/url"
	"os"
	"os/signal"
)

var (
	// ErrBuildMongoURIInvalidHost represents an error returned by the BuildMongoURI function when the invalid host
	// is used.
	ErrBuildMongoURIInvalidHost = errors.New("invalid host")
)

// GracefulShutdown creates a channel for signaling a graceful shutdown.
// When either an os.Interrupt or os.Kill signal is received, the returned channel is closed.
// This channel can be used by other parts of the program to wait for the shutdown signal.
func GracefulShutdown() <-chan struct{} {
	shutdown := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)
		signal.Notify(s, os.Interrupt, os.Kill)

		<-s
		close(shutdown)
	}()
	return shutdown
}

// BuildMongoURI builds a MongoDB connection URI from the provided host, port, username, and password.
// The returned URI string is of the format: "mongodb://[username:password@]host:port".
//
// If the port is not provided, the scheme of the URI is set to "mongodb+srv" instead of "mongodb".
// In addition, use of the +srv connection string modifier automatically sets the tls (or the equivalent ssl) option to
// true for the connection.
func BuildMongoURI(host, port, username, password string) (string, error) {
	if len(host) == 0 {
		return "", ErrBuildMongoURIInvalidHost
	}

	var userinfo *url.Userinfo
	if len(username) > 0 && len(password) > 0 {
		userinfo = url.UserPassword(username, password)
	} else if len(username) > 0 {
		userinfo = url.User(username)
	}

	proto := "mongodb"
	if len(port) > 0 {
		host = net.JoinHostPort(host, port)
	} else {
		proto += "+srv"
	}
	uri := url.URL{Scheme: proto, User: userinfo, Host: host}
	return uri.String(), nil
}
