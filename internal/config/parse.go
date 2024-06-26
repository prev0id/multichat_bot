package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func Parse(path string) (*Config, error) {
	bytes, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = json.Unmarshal(bytes, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
