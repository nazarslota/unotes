package mongo

import (
	"context"
	"net"
	"net/url"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// NewClient creates a new MongoDB client, ping and returns it.
func NewClient(ctx context.Context, uri string) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}

// BuildURI builds new MongoDB URI.
func BuildURI(host, port, username, password string) string {
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
	return uri.String()
}
