// Package mongo provides a MongoDB database repository implementation.
package mongo

import (
	"context"
	"fmt"

	"github.com/nazarslota/unotes/auth/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Config represents MongoDB configuration.
type Config struct {
	Host     string // Host specifies the MongoDB server host.
	Port     string // Port specifies the MongoDB server port.
	Username string // Username specifies the username used to authenticate with the MongoDB server.
	Password string // Password specifies the password used to authenticate with the MongoDB server.
	Database string // Database specifies the name of the MongoDB database to use.
}

// NewMongoDB returns a new instance of the *mongo.Database type using the provided configuration.
func NewMongoDB(ctx context.Context, config Config) (*mongo.Database, error) {
	uri, err := utils.BuildMongoURI(config.Host, config.Port, config.Username, config.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to build mongo uri: %w", err)
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to create a mongo client: %w", err)
	}

	if err := client.Connect(ctx); err != nil {
		return nil, fmt.Errorf("failed to establish connection with mono: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping mongo: %w", err)
	}
	return client.Database(config.Database), nil
}
