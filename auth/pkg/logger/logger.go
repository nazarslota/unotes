package logger

import (
	"fmt"
	"io"

	"github.com/rs/zerolog"
)

// Logger struct wraps a zerolog logger object and provides methods for logging at different levels.
type Logger struct {
	logger zerolog.Logger
}

// NewLogger creates and returns a new Logger.
// io.Writer is the destination to which the logs will be written.
func NewLogger(w io.Writer) Logger {
	l := zerolog.New(w)
	return Logger{logger: l}
}

// Trace logs a message at trace level.
func (l Logger) Trace(v ...any) {
	l.logger.Trace().Msg(fmt.Sprint(v...))
}

// Tracef logs a message at trace level with a format string.
func (l Logger) Tracef(format string, v ...any) {
	l.logger.Trace().Msgf(format, v...)
}

// TraceFields logs a message at trace level with additional fields.
func (l Logger) TraceFields(msg string, fields map[string]any) {
	l.logger.Trace().Fields(fields).Msg(msg)
}

// Debug logs a message at debug level.
func (l Logger) Debug(v ...any) {
	l.logger.Debug().Msg(fmt.Sprint(v...))
}

// Debugf logs a message at debug level with a format string.
func (l Logger) Debugf(format string, v ...any) {
	l.logger.Debug().Msgf(format, v...)
}

// DebugFields logs a message at debug level with additional fields.
func (l Logger) DebugFields(msg string, fields map[string]any) {
	l.logger.Debug().Fields(fields).Msg(msg)
}

// Info logs a message at info level.
func (l Logger) Info(v ...any) {
	l.logger.Info().Msg(fmt.Sprint(v...))
}

// Infof logs a message at info level with a format string.
func (l Logger) Infof(format string, v ...any) {
	l.logger.Info().Msgf(format, v...)
}

// InfoFields logs a message at info level with additional fields.
func (l Logger) InfoFields(msg string, fields map[string]any) {
	l.logger.Info().Fields(fields).Msg(msg)
}

// Warn method logs a warning message with the given variable arguments.
func (l Logger) Warn(v ...any) {
	l.logger.Warn().Msg(fmt.Sprint(v...))
}

// Warnf method logs a warning message with the given format string and variable arguments.
func (l Logger) Warnf(format string, v ...any) {
	l.logger.Warn().Msgf(format, v...)
}

// WarnFields method logs a warning message with the given fields.
func (l Logger) WarnFields(msg string, fields map[string]any) {
	l.logger.Warn().Fields(fields).Msg(msg)
}

// Error method logs an error message with the given variable arguments.
func (l Logger) Error(v ...any) {
	l.logger.Error().Msg(fmt.Sprint(v...))
}

// Errorf method logs an error message with the given format string and variable arguments.
func (l Logger) Errorf(format string, v ...any) {
	l.logger.Error().Msgf(format, v...)
}

// ErrorFields method logs an error message with the given fields.
func (l Logger) ErrorFields(msg string, fields map[string]any) {
	l.logger.Error().Fields(fields).Msg(msg)
}

// Fatal method logs a fatal message with the given variable arguments and causes the application to exit.
func (l Logger) Fatal(v ...any) {
	l.logger.Fatal().Msg(fmt.Sprint(v...))
}

// Fatalf method logs a fatal message with the given format string and variable arguments and causes the application to exit.
func (l Logger) Fatalf(format string, v ...any) {
	l.logger.Fatal().Msgf(format, v...)
}

// FatalFields method logs a fatal message with the given fields and causes the application to exit.
func (l Logger) FatalFields(msg string, fields map[string]any) {
	l.logger.Fatal().Fields(fields).Msg(msg)
}

// Panic method logs a panic message with the given variable arguments and causes a panic.
func (l Logger) Panic(v ...any) {
	l.logger.Panic().Msg(fmt.Sprint(v...))
}

// Panicf method logs a panic message with the given format string and variable arguments and causes a panic.
func (l Logger) Panicf(format string, v ...any) {
	l.logger.Panic().Msgf(format, v...)
}

// PanicFields method logs a panic message with the given fields and causes a panic.
func (l Logger) PanicFields(msg string, fields map[string]any) {
	l.logger.Panic().Fields(fields).Msg(msg)
}

// Context struct holds a zerolog.Logger object and provides methods for creating new contexts with
// additional fields.
type Context struct {
	logger zerolog.Logger
}

// With returns a new context with the same logger.
func (l Logger) With() Context {
	return Context{logger: l.logger}
}

// Logger returns the logger from the context.
func (c Context) Logger() Logger {
	return Logger{logger: c.logger}
}

// Timestamp adds a timestamp field to the context's logger.
// Returns a new context with the updated logger.
func (c Context) Timestamp() Context {
	return Context{logger: c.logger.With().Timestamp().Logger()}
}

// Caller adds a caller field to the context's logger.
// Returns a new context with the updated logger.
func (c Context) Caller() Context {
	return Context{logger: c.logger.With().CallerWithSkipFrameCount(3).Logger()}
}
