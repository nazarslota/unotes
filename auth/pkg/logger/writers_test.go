package logger

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConsoleWriter(t *testing.T) {
	buf := new(bytes.Buffer)
	writer := NewConsoleWriter(func(w *ConsoleWriter) {
		w.Out = buf
		w.NoColor = false
		w.TimeFormat = time.RFC1123
	})
	require.NotNil(t, writer)

	NewLogger(writer).Info("Info message")
	assert.Equal(t, "\x1b[90m<nil>\x1b[0m \x1b[32mINF\x1b[0m Info message\n", buf.String())
}
