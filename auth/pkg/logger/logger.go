// Package logger provides Logger, which is a wrapper around zerolog.Logger.
package logger

import (
	"fmt"
	"io"

	"github.com/rs/zerolog"
)

// Logger provides a structured and leveled logging interface using zerolog.Logger.
type Logger struct {
	logger zerolog.Logger
}

// NewLogger creates a new Logger instance that writes to the specified io.Writer.
func NewLogger(w io.Writer) Logger {
	l := zerolog.New(w)
	return Logger{logger: l}
}

// Trace logs a message at the Trace severity level.
func (l Logger) Trace(v ...any) {
	l.logger.Trace().Msg(fmt.Sprint(v...))
}

// Tracef logs a message at the Trace severity level.
func (l Logger) Tracef(format string, v ...any) {
	l.logger.Trace().Msgf(format, v...)
}

// TraceFields logs a message with associated fields at the Trace severity level.
func (l Logger) TraceFields(msg string, fields map[string]any) {
	l.logger.Trace().Fields(fields).Msg(msg)
}

// Debug logs a message at the Debug severity level.
func (l Logger) Debug(v ...any) {
	l.logger.Debug().Msg(fmt.Sprint(v...))
}

// Debugf logs a message at the Debug severity level.
func (l Logger) Debugf(format string, v ...any) {
	l.logger.Debug().Msgf(format, v...)
}

// DebugFields logs a message with associated fields at the Debug severity level.
func (l Logger) DebugFields(msg string, fields map[string]any) {
	l.logger.Debug().Fields(fields).Msg(msg)
}

// Info logs a message at the Info severity level.
func (l Logger) Info(v ...any) {
	l.logger.Info().Msg(fmt.Sprint(v...))
}

// Infof logs a message at the Info severity level.
func (l Logger) Infof(format string, v ...any) {
	l.logger.Info().Msgf(format, v...)
}

// InfoFields logs a message with associated fields at the Info severity level.
func (l Logger) InfoFields(msg string, fields map[string]any) {
	l.logger.Info().Fields(fields).Msg(msg)
}

// Warn logs a message at the Warn severity level.
func (l Logger) Warn(v ...any) {
	l.logger.Warn().Msg(fmt.Sprint(v...))
}

// Warnf logs a message at the Warn severity level.
func (l Logger) Warnf(format string, v ...any) {
	l.logger.Warn().Msgf(format, v...)
}

// WarnFields logs a message with associated fields at the Warn severity level.
func (l Logger) WarnFields(msg string, fields map[string]any) {
	l.logger.Warn().Fields(fields).Msg(msg)
}

// Error logs a message at the Error severity level.
func (l Logger) Error(v ...any) {
	l.logger.Error().Msg(fmt.Sprint(v...))
}

// Errorf logs a message at the Error severity level.
func (l Logger) Errorf(format string, v ...any) {
	l.logger.Error().Msgf(format, v...)
}

// ErrorFields logs a message with associated fields at the Error severity level.
func (l Logger) ErrorFields(msg string, fields map[string]any) {
	l.logger.Error().Fields(fields).Msg(msg)
}

// Fatal logs a message at the Fatal severity level.
func (l Logger) Fatal(v ...any) {
	l.logger.Fatal().Msg(fmt.Sprint(v...))
}

// Fatalf logs a message at the Fatal severity level.
func (l Logger) Fatalf(format string, v ...any) {
	l.logger.Fatal().Msgf(format, v...)
}

// FatalFields logs a message with associated fields at the Fatal severity level.
func (l Logger) FatalFields(msg string, fields map[string]any) {
	l.logger.Fatal().Fields(fields).Msg(msg)
}

// Panic logs a message at the Panic severity level.
func (l Logger) Panic(v ...any) {
	l.logger.Panic().Msg(fmt.Sprint(v...))
}

// Panicf logs a message at the Panic severity level.
func (l Logger) Panicf(format string, v ...any) {
	l.logger.Panic().Msgf(format, v...)
}

// PanicFields logs a message with associated fields at the Panic severity level.
func (l Logger) PanicFields(msg string, fields map[string]any) {
	l.logger.Panic().Fields(fields).Msg(msg)
}

// Context configures a new sub-logger with contextual fields.
type Context struct {
	logger zerolog.Logger
}

// With creates a child logger with the field added to its context.
func (l Logger) With() Context {
	return Context{logger: l.logger}
}

// Timestamp adds the current local time as UNIX timestamp to the logger context with the "time" key.
// To customize the key name, change zerolog.TimestampFieldName.
func (c Context) Timestamp() Context {
	return Context{logger: c.logger.With().Timestamp().Logger()}
}

// Caller adds the file:line of the caller with the zerolog.CallerFieldName key.
func (c Context) Caller() Context {
	return Context{logger: c.logger.With().CallerWithSkipFrameCount(3).Logger()}
}

// Logger returns the logger with the context previously set.
func (c Context) Logger() Logger {
	return Logger{logger: c.logger}
}
