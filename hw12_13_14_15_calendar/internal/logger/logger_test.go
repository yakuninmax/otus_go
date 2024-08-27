package logger

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/config"
)

func TestLogger(t *testing.T) {
	// Test cases.
	testcases := []struct {
		Name   string
		Config config.LoggerConfig
		Result string
		Action func(logger *Logger)
	}{
		{
			Name: "no message",
			Config: config.LoggerConfig{
				Level:         "info",
				Colors:        false,
				FullTimestamp: false,
			},
			Result: "",
			Action: func(logger *Logger) {
				logger.Debug("Debug message")
			},
		},
		{
			Name: "info message",
			Config: config.LoggerConfig{
				Level:         "info",
				Colors:        false,
				FullTimestamp: false,
			},
			Result: "Info message",
			Action: func(logger *Logger) {
				logger.Info("Info message")
			},
		},
		{
			Name: "warning message",
			Config: config.LoggerConfig{
				Level:         "warn",
				Colors:        false,
				FullTimestamp: false,
			},
			Result: "Warning message",
			Action: func(logger *Logger) {
				logger.Warn("Warning message")
			},
		},
		{
			Name: "error message",
			Config: config.LoggerConfig{
				Level:         "error",
				Colors:        false,
				FullTimestamp: false,
			},
			Result: "Error message",
			Action: func(logger *Logger) {
				logger.Error("Error message")
			},
		},
		{
			Name: "debug message",
			Config: config.LoggerConfig{
				Level:         "debug",
				Colors:        false,
				FullTimestamp: false,
			},
			Result: "Debug message",
			Action: func(logger *Logger) {
				logger.Debug("Debug message")
			},
		},
	}

	t.Run("message tests", func(t *testing.T) {
		for _, testcase := range testcases {
			stdoutBackup := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			logger, err := New(&testcase.Config)
			testcase.Action(logger)
			require.NoError(t, err)
			require.NotNil(t, logger)

			outC := make(chan string)
			go func() {
				var buf bytes.Buffer
				io.Copy(&buf, r)
				outC <- buf.String()
			}()

			w.Close()
			os.Stdout = stdoutBackup
			output := <-outC

			require.Contains(t, output, testcase.Result)
		}
	})

	t.Run("unsupported logging level", func(t *testing.T) {
		logger, err := New(&config.LoggerConfig{
			Level:         "unknown",
			Colors:        false,
			FullTimestamp: false,
		})
		require.Nil(t, logger)
		require.Error(t, err, ErrUnsupportedLoggingLevel)
	})
}
