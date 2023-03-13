package mongo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMongoDB(t *testing.T) {
	canceled, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		name     string
		ctx      context.Context
		host     string
		port     string
		username string
		password string
		database string
		wantErr  bool
	}{
		{
			name:     "should return new mongo database client",
			ctx:      context.Background(),
			host:     "localhost",
			port:     "27017",
			username: "",
			password: "",
			wantErr:  false,
		},
		{
			name:     "should return an error if host is invalid",
			ctx:      context.Background(),
			host:     "invalid-host",
			port:     "27017",
			username: "",
			password: "",
			wantErr:  true,
		},
		{
			name:     "should return an error if port is invalid",
			ctx:      context.Background(),
			host:     "localhost",
			port:     "-1",
			username: "",
			password: "",
			wantErr:  true,
		},
		{
			name:     "should return an error if context is invalid",
			ctx:      canceled,
			host:     "localhost",
			port:     "27017",
			username: "",
			password: "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := NewMongoDB(tt.ctx, Config{
				Host:     tt.host,
				Port:     tt.port,
				Username: tt.username,
				Password: tt.password,
				Database: tt.database,
			})

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)
			}
		})
	}
}
