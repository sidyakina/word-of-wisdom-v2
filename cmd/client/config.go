package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"time"
)

type Config struct {
	ServerURL       string        `env:"SERVER_URL" envDefault:"localhost:8080"`
	ReadTimeout     time.Duration `env:"READ_TIMEOUT" envDefault:"2s"`
	SolutionTimeout time.Duration `env:"SOLUTION_TIMEOUT" envDefault:"30s"`
}

func parseConfig() (*Config, error) {
	cfg := &Config{}

	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}

	if cfg.ServerURL == "" {
		return nil, fmt.Errorf("empty server URL")
	}

	if cfg.ReadTimeout <= 0 {
		return nil, fmt.Errorf("wrong read timeout %v", cfg.ReadTimeout)
	}

	if cfg.SolutionTimeout <= 0 {
		return nil, fmt.Errorf("wrong solution timeout %v", cfg.ReadTimeout)
	}

	return cfg, nil
}
