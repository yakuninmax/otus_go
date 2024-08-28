package config

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestConfig(t *testing.T) {
	// Reference config structure.
	referenceConfig := Config{
		LoggerConfig: LoggerConfig{
			Level:         "debug",
			Colors:        true,
			FullTimestamp: true,
		},
		StorageConfig: StorageConfig{
			Type: "in-memory",
			StorageConnection: StorageConnection{
				Database: "calendar",
				User:     "calendar",
				Password: "c@1end@r",
				Host:     "db",
				Port:     5432,
			},
		},
		HTTPConfig: HTTPConfig{
			Host:    "localhost",
			Port:    8080,
			LogPath: "./calendar-http.log",
		},
	}

	t.Run("compare configs", func(t *testing.T) {
		actualConfig, err := NewConfig("../../configs/config.yaml")
		require.NoError(t, err)
		require.Equal(t, &referenceConfig, actualConfig)
	})

	t.Run("invalid config file", func(t *testing.T) {
		_, err := NewConfig("../../build/Dockerfile")
		require.Error(t, err, &yaml.TypeError{})
	})
}
