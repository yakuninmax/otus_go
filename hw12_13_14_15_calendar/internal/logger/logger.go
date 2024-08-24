package logger

import (
	"errors"
	"os"

	"github.com/sirupsen/logrus"
)

var ErrUnsupportedLoggingLevel = errors.New("unsupported logging level")

// Logger structure.
type Logger struct {
	logger *logrus.Logger
}

// New logger.
func New(level string, colors, fullTimestamp bool) (*Logger, error) {
	// Create new logger.
	logger := logrus.New()

	// Setup text formatter.
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: colors,
		FullTimestamp: fullTimestamp,
	})

	// Set output channel.
	logger.SetOutput(os.Stdout)

	// Set logging level.
	// Acceptable values: error, warn, info, debug.
	switch level {
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	default:
		return nil, ErrUnsupportedLoggingLevel
	}

	return &Logger{logger}, nil
}

func (l Logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l Logger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l Logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l Logger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}
