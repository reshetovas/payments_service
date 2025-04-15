package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	LogLevel  string
	CacheType string
}

// Is it nevessary to do the method
func LoadConfig() (*Config, error) {
	file, err := os.ReadFile("config.json")
	if err != nil {
		return nil, fmt.Errorf("error with reading config.json: %w", err)
	}

	var cfg Config
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return nil, fmt.Errorf("config Json parsing error: %w", err)
	}

	return &cfg, nil
}
