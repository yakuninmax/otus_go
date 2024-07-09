package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("Successful execution", func(t *testing.T) {
		result := RunCmd([]string{"uname", "-a"}, nil)
		require.Equal(t, 0, result)
	})

	t.Run("Failed execution", func(t *testing.T) {
		result := RunCmd([]string{"cp", "a", "b"}, nil)
		require.Equal(t, 1, result)
	})
}
