package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"time"
)

type Config struct {
	Port uint32 `env:"PORT" envDefault:"8080"`
	// read timeout is big because client need time to generate solution
	ReadSolutionTimeout time.Duration `env:"READ_SOLUTION_TIMEOUT" envDefault:"15s"`

	LenChallengeString int32 `env:"LEN_CHALLENGE_STRING" envDefault:"20"`
	LenSolutionString  int32 `env:"LEN_SOLUTION_STRING" envDefault:"10"`
	NumberZeroBits     int32 `env:"NUMBER_ZERO_BITS" envDefault:"20"`
}

func parseConfig() (*Config, error) {
	cfg := &Config{}

	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}

	if cfg.Port == 0 {
		return nil, fmt.Errorf("zero server port")
	}

	if cfg.ReadSolutionTimeout <= 0 {
		return nil, fmt.Errorf("wrong ReadSolutionTimeout %v", cfg.ReadSolutionTimeout)
	}

	if cfg.LenChallengeString <= 0 {
		return nil, fmt.Errorf("wrong LenChallengeString %v", cfg.LenChallengeString)
	}

	if cfg.LenSolutionString <= 0 {
		return nil, fmt.Errorf("wrong LenSolutionString %v", cfg.LenSolutionString)
	}

	if cfg.NumberZeroBits <= 0 {
		return nil, fmt.Errorf("wrong NumberZeroBits %v", cfg.NumberZeroBits)
	}

	return cfg, nil
}
