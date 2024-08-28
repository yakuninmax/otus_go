package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Config struct.
type Config struct {
	LoggerConfig  `yaml:"log"`
	StorageConfig `yaml:"storage"`
	HTTPConfig    `yaml:"http"`
}

// Logger config structure.
type LoggerConfig struct {
	// Log level.
	// Acceptable values: error, warn, info, debug.
	Level         string `yaml:"level"`
	Colors        bool   `yaml:"colors"`
	FullTimestamp bool   `yaml:"fullTimestamp"`
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

	// Log file path.
	LogPath string `yaml:"logPath"`
}

func NewConfig(configFile string) (*Config, error) {
	// Check parameters.
	fileInfo, err := os.Stat(configFile)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	if fileInfo.IsDir() {
		return nil, fmt.Errorf("'%s' is a directory", configFile)
	}

	// Open config file.
	file, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %w", err)
	}
	defer file.Close()

	// YAML decoding from file.
	config := Config{}

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
