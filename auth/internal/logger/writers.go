package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type TextWriter struct {
	Out        io.Writer
	NoColor    bool
	TimeFormat string

	console *zerolog.ConsoleWriter
}

func NewTextWriter(options ...func(w *TextWriter)) *TextWriter {
	w := &TextWriter{
		Out:        os.Stdout,
		NoColor:    false,
		TimeFormat: time.RFC822,
	}

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

func (w *TextWriter) Write(p []byte) (n int, err error) {
	if w.console != nil {
		return w.console.Write(p)
	}

	w.console = &zerolog.ConsoleWriter{
		Out:        w.Out,
		NoColor:    w.NoColor,
		TimeFormat: w.TimeFormat,
	}
	return w.console.Write(p)
}
