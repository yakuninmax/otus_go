package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config struct.
type Config struct {
	LoggerConfig  `yaml:",inline"`
	StorageConfig `yaml:"storage"`
	HTTPConfig    `yaml:"http"`
}

// Logger config structure.
type LoggerConfig struct {
	// Log level.
	// Acceptable values: error, warn, info, debug.
	Level string `yaml:"logLevel"`
}

// Storage parameters structure.
type StorageConfig struct {
	// Storage type.
	// Acceptable values: in-memory, postgresql.
	Type              string `yaml:"type"`
	StorageConnection `yaml:"connection"`
}

// Storage connection structure.
type StorageConnection struct {
	// Database name.
	Database string `yaml:"database"`

	// User name.
	User string `yaml:"user"`

	// User password.
	Password string `yaml:"password"`

	// Database server host.
	Host string `yaml:"host"`

	// Database server port.
	Port int `yaml:"port"`
}

// HTTP server parameters.
type HTTPConfig struct {
	// HTTP server host.
	Host string `yaml:"host"`

	// HTTP server port.
	Port int `yaml:"port"`
}

func NewConfig(configFile string) (Config, error) {
	// Create config structure.
	config := Config{}

	// Open config file.
	file, err := os.Open(configFile)
	if err != nil {
		return config, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}
