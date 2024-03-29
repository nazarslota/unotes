package utils

import (
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGracefulShutdown(t *testing.T) {
	shutdown := GracefulShutdown()
	assert.NotNil(t, shutdown)

	go func() {
		time.Sleep(100 * time.Millisecond)

		err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		require.NoError(t, err)
	}()

	select {
	case <-time.After(15 * time.Second):
		require.Fail(t, "Timed out waiting for shutdown")
	case <-shutdown:
	}
}

func TestBuildMongoURI(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		port     string
		username string
		password string
		expected string
		err      error
	}{
		{
			name:     "default MongoDB URI",
			host:     "localhost",
			port:     "27017",
			username: "",
			password: "",
			expected: "mongodb://localhost:27017",
			err:      nil,
		},
		{
			name:     "MongoDB URI with username and password",
			host:     "cluster0.mongodb.net",
			port:     "",
			username: "user",
			password: "pass",
			expected: "mongodb+srv://user:pass@cluster0.mongodb.net",
			err:      nil,
		},
		{
			name:     "MongoDB URI with only username",
			host:     "localhost",
			port:     "",
			username: "user",
			password: "",
			expected: "mongodb+srv://user@localhost",
			err:      nil,
		},
		{
			name:     "invalid MongoDB URI",
			host:     "",
			port:     "",
			username: "",
			password: "",
			expected: "",
			err:      ErrBuildMongoURIInvalidHost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uri, err := BuildMongoURI(tt.host, tt.port, tt.username, tt.password)
			assert.Equal(t, tt.expected, uri)
			assert.ErrorIs(t, err, tt.err)
		})
	}
}
