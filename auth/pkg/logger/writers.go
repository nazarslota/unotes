package logger

import (
	"io"

	"github.com/rs/zerolog"
)

// ConsoleWriter is a struct that wraps the zerolog ConsoleWriter with additional options.
type ConsoleWriter struct {
	// Out is the io.Writer to write logs to.
	Out io.Writer
	// NoColor controls whether the console output should be colorized.
	NoColor bool
	// TimeFormat is the format of the timestamp in the log message.
	TimeFormat string

	// console is the underlying zerolog ConsoleWriter
	console *zerolog.ConsoleWriter
}

// NewTextWriter creates a new ConsoleWriter with the given options applied.
func NewTextWriter(options ...func(w *ConsoleWriter)) *ConsoleWriter {
	w := &ConsoleWriter{}
	for _, option := range options {
		option(w)
	}

	w.console = &zerolog.ConsoleWriter{
		Out:        w.Out,
		NoColor:    w.NoColor,
		TimeFormat: w.TimeFormat,
	}
	return w
}

// Write implements the io.Writer interface and makes calls to the base zerolog ConsoleWriter.
func (w ConsoleWriter) Write(p []byte) (n int, err error) {
	if w.console == nil {
		w.console = &zerolog.ConsoleWriter{
			Out:        w.Out,
			NoColor:    w.NoColor,
			TimeFormat: w.TimeFormat,
		}
	}
	return w.console.Write(p)
}
