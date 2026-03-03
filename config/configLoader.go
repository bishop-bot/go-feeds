package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Project    string   `yaml:"project"`
	MaxWorkers int      `yaml:"maxWorkers"`
	Feeds      []string `yaml:"feeds"`
}

func NewConfig() Config {
	return Config{
		Project:    "go-feeds",
		MaxWorkers: 5,
		Feeds:      []string{},
	}
}

// LoadConfig loads the configuration from a YAML file.
func LoadConfig(filename string) (*Config, error) {
	// default config
	config := NewConfig()

	// Read the YAML file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal the YAML data into the Config struct
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// MustLoadConfig loads the configuration and panics if there is an error.
func MustLoadConfig(filename string) *Config {
	config, err := LoadConfig(filename)
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading config")
	}
	return config
}
