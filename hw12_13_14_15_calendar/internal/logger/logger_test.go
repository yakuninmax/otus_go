package logger

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	// Test cases.
	testcases := []struct {
		Name   string
		Level  string
		Result string
		Action func(logger *Logger)
	}{
		{
			Name:   "no message",
			Level:  "info",
			Result: "",
			Action: func(logger *Logger) {
				logger.Debug("Debug message")
			},
		},
		{
			Name:   "info message",
			Level:  "info",
			Result: "Info message",
			Action: func(logger *Logger) {
				logger.Info("Info message")
			},
		},
		{
			Name:   "warning message",
			Level:  "warn",
			Result: "Warning message",
			Action: func(logger *Logger) {
				logger.Warn("Warning message")
			},
		},
		{
			Name:   "error message",
			Level:  "error",
			Result: "Error message",
			Action: func(logger *Logger) {
				logger.Error("Error message")
			},
		},
		{
			Name:   "debug message",
			Level:  "debug",
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

			logger, err := New(testcase.Level, true, true)
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
		logger, err := New("trace", true, true)
		require.Nil(t, logger)
		require.Error(t, err, ErrUnsupportedLoggingLevel)
	})
}
