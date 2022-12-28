package mongo

import (
	"context"
	"fmt"

	"github.com/udholdenhed/unotes/auth/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func NewMongoDB(ctx context.Context, config *Config) (*mongo.Database, error) {
	uri := utils.BuildMongoURI(
		config.Host,
		config.Port,
		config.Username,
		config.Password,
	)

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
