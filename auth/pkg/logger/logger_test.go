package logger

import (
	"bytes"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	buf := new(bytes.Buffer)
	fields := map[string]any{"key1": "value1", "key2": "value2"}

	logger := NewLogger(buf)
	require.NotNil(t, logger)

	t.Run("trace level", func(t *testing.T) {
		logger.Trace("This is a Trace log message")
		assert.Contains(t, buf.String(), "This is a Trace log message")

		logger.Tracef("This is a formatted Trace log message with a %s", "parameter")
		assert.Contains(t, buf.String(), "This is a formatted Trace log message with a parameter")

		logger.TraceFields("This is a Trace log message with fields", fields)
		assert.Contains(t, buf.String(), "This is a Trace log message with fields")
		assert.Contains(t, buf.String(), `"key1":"value1"`)
		assert.Contains(t, buf.String(), `"key2":"value2"`)

		t.Cleanup(func() { buf.Reset() })
	})

	t.Run("debug level", func(t *testing.T) {
		logger.Debug("This is a Debug log message")
		assert.Contains(t, buf.String(), "This is a Debug log message")

		logger.Debugf("This is a formatted Debug log message with a %s", "parameter")
		assert.Contains(t, buf.String(), "This is a formatted Debug log message with a parameter")

		logger.DebugFields("This is a Debug log message with fields", fields)
		assert.Contains(t, buf.String(), "This is a Debug log message with fields")
		assert.Contains(t, buf.String(), `"key1":"value1"`)
		assert.Contains(t, buf.String(), `"key2":"value2"`)

		t.Cleanup(func() { buf.Reset() })
	})

	t.Run("info level", func(t *testing.T) {
		logger.Info("This is a Info log message")
		assert.Contains(t, buf.String(), "This is a Info log message")

		logger.Infof("This is a formatted Info log message with a %s", "parameter")
		assert.Contains(t, buf.String(), "This is a formatted Info log message with a parameter")

		logger.InfoFields("This is a Info log message with fields", fields)
		assert.Contains(t, buf.String(), "This is a Info log message with fields")
		assert.Contains(t, buf.String(), `"key1":"value1"`)
		assert.Contains(t, buf.String(), `"key2":"value2"`)

		t.Cleanup(func() { buf.Reset() })
	})

	t.Run("warn level", func(t *testing.T) {
		logger.Warn("This is a Warn log message")
		assert.Contains(t, buf.String(), "This is a Warn log message")

		logger.Warnf("This is a formatted Warn log message with a %s", "parameter")
		assert.Contains(t, buf.String(), "This is a formatted Warn log message with a parameter")

		logger.WarnFields("This is a Warn log message with fields", fields)
		assert.Contains(t, buf.String(), "This is a Warn log message with fields")
		assert.Contains(t, buf.String(), `"key1":"value1"`)
		assert.Contains(t, buf.String(), `"key2":"value2"`)

		t.Cleanup(func() { buf.Reset() })
	})

	t.Run("error level", func(t *testing.T) {
		logger.Error("This is a Error log message")
		assert.Contains(t, buf.String(), "This is a Error log message")

		logger.Errorf("This is a formatted Error log message with a %s", "parameter")
		assert.Contains(t, buf.String(), "This is a formatted Error log message with a parameter")

		logger.ErrorFields("This is a Error log message with fields", fields)
		assert.Contains(t, buf.String(), "This is a Error log message with fields")
		assert.Contains(t, buf.String(), `"key1":"value1"`)
		assert.Contains(t, buf.String(), `"key2":"value2"`)

		t.Cleanup(func() { buf.Reset() })
	})

	t.Run("fatal level", func(t *testing.T) {
		t.Run("fatal", func(t *testing.T) {
			if os.Getenv("FLAG") == "1" {
				logger.Fatal("This is a Fatal log message")
				return
			}

			cmd := exec.Command(os.Args[0], "-test.run=TestLogger/fatal_level/fatal")
			cmd.Env = append(cmd.Env, "FLAG=1")

			err, ok := cmd.Run().(*exec.ExitError)
			assert.True(t, ok)
			assert.Equal(t, 1, err.ExitCode())
		})

		t.Run("fatalf", func(t *testing.T) {
			if os.Getenv("FLAG") == "1" {
				logger.Fatalf("This is a formatted Fatal log message with a %s", "parameter")
				return
			}

			cmd := exec.Command(os.Args[0], "-test.run=TestLogger/fatal_level/fatalf")
			cmd.Env = append(cmd.Env, "FLAG=1")

			err, ok := cmd.Run().(*exec.ExitError)
			assert.True(t, ok)
			assert.Equal(t, 1, err.ExitCode())
		})

		t.Run("fatal fields", func(t *testing.T) {
			if os.Getenv("FLAG") == "1" {
				logger.FatalFields("This is a Fatal log message with fields", fields)
				return
			}

			cmd := exec.Command(os.Args[0], "-test.run=TestLogger/fatal_level/fatal_fields")
			cmd.Env = append(cmd.Env, "FLAG=1")

			err, ok := cmd.Run().(*exec.ExitError)
			assert.True(t, ok)
			assert.Equal(t, 1, err.ExitCode())
		})

		t.Cleanup(func() { buf.Reset() })
	})

	t.Run("panic level", func(t *testing.T) {
		t.Run("panic", func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					assert.Contains(t, buf.String(), "This is a Panic log message")
				}
			}()
			logger.Panic("This is a Panic log message")
		})

		t.Run("panicf", func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					assert.Contains(t, buf.String(), "This is a formatted Panic log message with a parameter")
				}
			}()
			logger.Panicf("This is a formatted Panic log message with a %s", "parameter")
		})

		t.Run("panic fields", func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					assert.Contains(t, buf.String(), "This is a Panic log message with fields")
					assert.Contains(t, buf.String(), `"key1":"value1"`)
					assert.Contains(t, buf.String(), `"key2":"value2"`)
				}
			}()
			logger.PanicFields("This is a Panic log message with fields", fields)
		})

		t.Cleanup(func() { buf.Reset() })
	})
}
