package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	// Reference file name.
	envDir := "./testdata/env"

	// Create reference map.
	refMap := make(Environment)
	refMap["BAR"] = EnvValue{Value: "bar"}
	refMap["EMPTY"] = EnvValue{Value: ""}
	refMap["FOO"] = EnvValue{Value: "   foo\nwith new line"}
	refMap["HELLO"] = EnvValue{Value: "\"hello\""}
	refMap["UNSET"] = EnvValue{NeedRemove: true}

	t.Run("Read envs from folder", func(t *testing.T) {
		envVars, err := ReadDir(envDir)
		require.NoError(t, err)
		require.Equal(t, refMap, envVars)
	})

	t.Run("Test equal sign in file name", func(t *testing.T) {
		// Create file with equal sing in the name.
		fileName := path.Join(envDir, "A=K")
		_, err := os.Create(fileName)

		require.NoError(t, err)

		defer os.Remove(fileName)

		_, err = ReadDir(envDir)
		require.ErrorIs(t, err, errInvalidFileName)
	})
}
