package logger

import (
	"fmt"
	"io"

	"github.com/rs/zerolog"
)

type Logger struct {
	logger *zerolog.Logger
}

func New(w io.Writer) *Logger {
	l := zerolog.New(w)
	return &Logger{&l}
}

func (l *Logger) Info(v ...any) {
	l.logger.Info().Msg(fmt.Sprint(v...))
}

func (l *Logger) Warn(v ...any) {
	l.logger.Warn().Msg(fmt.Sprint(v...))
}

func (l *Logger) Error(args ...any) {
	l.logger.Error().Msg(fmt.Sprint(args...))
}

func (l *Logger) Infof(format string, v ...any) {
	l.logger.Info().Msgf(format, v...)
}

func (l *Logger) Warnf(format string, v ...any) {
	l.logger.Warn().Msgf(format, fmt.Sprint(v...))
}

func (l *Logger) Errorf(format string, args ...any) {
	l.logger.Error().Msgf(format, fmt.Sprint(args...))
}

func (l *Logger) InfoFields(msg string, fields map[string]any) {
	l.logger.Info().Fields(fields).Msg(msg)
}

func (l *Logger) WarnFields(msg string, fields map[string]any) {
	l.logger.Warn().Fields(fields).Msg(msg)
}

func (l *Logger) ErrorFields(msg string, fields map[string]any) {
	l.logger.Error().Fields(fields).Msg(msg)
}
