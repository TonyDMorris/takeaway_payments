package config

import "github.com/caarlos0/env/v6"

type Config struct {
	StrapiURL         string `env:"STRAPI_URL"`
	ServiceIdentifier string `env:"SERVICE_IDENTIFIER"`
	ServicePassword   string `env:"SERVICE_PASSWORD"`

	DBHost     string `env:"DB_HOST"`
	DBPort     int    `env:"DB_PORT"`
	DBUsername string `env:"DB_USERNAME"`
	DBPassword string `env:"DB_PASSWORD"`
	DBName     string `env:"DB_NAME"`
}

func GetConfig() *Config {
	cfg := &Config{}
	env.Parse(cfg)
	return cfg
}
