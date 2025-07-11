package config

import (
	"encoding/json"
	"os"
)

type (
	Config struct {
		LoadBalancer LoadBalancer `json:"load_balancer"`
		HTTP         HTTP         `json:"http"`
		Log          Log          `json:"log"`
		PG           PG           `json:"pg"`
	}

	LoadBalancer struct {
		Backends []string `json:"backends"`
	}

	HTTP struct {
		Port string `json:"http_port"`
	}

	Log struct {
		Level string `json:"log_level"`
	}

	PG struct {
		URL string `json:"url"`
	}
)

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
