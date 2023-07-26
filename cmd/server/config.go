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

	ChallengeNumberSymbols int32 `env:"CHALLENGE_NUMBER_SYMBOLS" envDefault:"20"`
	SolutionNumberSymbols  int32 `env:"SOLUTION_NUMBER_SYMBOLS" envDefault:"10"`
	NumberZeroBits         int32 `env:"NUMBER_ZERO_BITS" envDefault:"20"`
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

	if cfg.ChallengeNumberSymbols <= 0 {
		return nil, fmt.Errorf("wrong ChallengeNumberSymbols %v", cfg.ChallengeNumberSymbols)
	}

	if cfg.SolutionNumberSymbols <= 0 {
		return nil, fmt.Errorf("wrong SolutionNumberSymbols %v", cfg.SolutionNumberSymbols)
	}

	if cfg.NumberZeroBits <= 0 {
		return nil, fmt.Errorf("wrong NumberZeroBits %v", cfg.NumberZeroBits)
	}

	return cfg, nil
}
