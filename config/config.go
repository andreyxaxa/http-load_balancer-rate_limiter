package config

import (
	"encoding/json"
	"os"
)

type (
	Config struct {
		LoadBalancer LoadBalancer
	}

	LoadBalancer struct {
		Backends []string `json:"backends"`
	}
)

// TODO: подумать, как сделать конфиг лучше, так читать неприятно
func NewConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
