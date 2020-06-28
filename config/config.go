package config

import "github.com/caarlos0/env/v6"

type Config struct {
	Something string `json:"something"`
}

func GetConfig() *Config {
	cfg := &Config{}
	env.Parse(cfg)
	return cfg
}
