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
		host     string
		port     string
		username string
		password string
		expected string
		err      error
	}{
		{
			host:     "localhost",
			port:     "27017",
			username: "",
			password: "",
			expected: "mongodb://localhost:27017",
			err:      nil,
		},
		{
			host:     "cluster0.mongodb.net",
			port:     "",
			username: "user",
			password: "pass",
			expected: "mongodb+srv://user:pass@cluster0.mongodb.net",
			err:      nil,
		},
		{
			host:     "localhost",
			port:     "",
			username: "user",
			password: "",
			expected: "mongodb+srv://user@localhost",
			err:      nil,
		},
		{
			host:     "",
			port:     "",
			username: "",
			password: "",
			expected: "",
			err:      ErrBuildMongoURIInvalidHost,
		},
	}

	for _, test := range tests {
		uri, err := BuildMongoURI(test.host, test.port, test.username, test.password)
		assert.Equal(t, test.err, err)
		assert.Equal(t, test.expected, uri)
	}
}
