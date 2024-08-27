package logger

import (
	"errors"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/config"
)

var ErrUnsupportedLoggingLevel = errors.New("unsupported logging level")

// Logger structure.
type Logger struct {
	logger *logrus.Logger
}

// New logger.
func New(config *config.LoggerConfig) (*Logger, error) {
	// Create new logger.
	logger := logrus.New()

	// Setup text formatter.
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: config.Colors,
		FullTimestamp: config.FullTimestamp,
	})

	// Set output channel.
	logger.SetOutput(os.Stdout)

	// Set logging level.
	// Acceptable values: error, warn, info, debug.
	switch config.Level {
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
