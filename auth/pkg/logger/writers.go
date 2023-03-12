package logger

import (
	"io"

	"github.com/rs/zerolog"
)

// ConsoleWriter is a console writer for Logger that writes log messages to the console in a user-friendly format.
type ConsoleWriter struct {
	Out        io.Writer // Out is the destination writer for log messages.
	NoColor    bool      // NoColor disables ANSI color escape codes in output.
	TimeFormat string    // TimeFormat is the format for timestamps in output.

	console *zerolog.ConsoleWriter // console is a ZeroLog ConsoleWriter used for writing output to the console.
}

// NewConsoleWriter creates a new ConsoleWriter with the specified options.
func NewConsoleWriter(options ...func(w *ConsoleWriter)) *ConsoleWriter {
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

// Write implements the io.Writer interface.
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
