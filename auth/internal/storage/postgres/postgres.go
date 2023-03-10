// Package postgres provides a PostgreSQL repository implementation for storing and managing users.
package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Config stores the configuration information required to connect to a PostgreSQL database.
type Config struct {
	Host     string // Host is the address of the PostgreSQL server.
	Port     string // Port is the port number to connect to on the PostgreSQL server.
	Username string // Username is the username used to authenticate against the PostgreSQL server.
	Password string // Password is the password used to authenticate against the PostgreSQL server.
	DBName   string // DBName is the name of the PostgreSQL database to connect to.
	SSLMode  string // SSLMode controls whether or with what priority a secure SSL TCP/IP connection will be negotiated with the PostgreSQL server.
}

// NewPostgreSQL creates a new PostgreSQL database connection using the provided configuration.
func NewPostgreSQL(ctx context.Context, config Config) (*sqlx.DB, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("invalid context: %w", ctx.Err())
	default:
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.DBName, config.SSLMode)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
