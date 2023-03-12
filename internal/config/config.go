package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Server
	Database
	Auth
}

type Auth struct {
	SecretKey string `env:"SECRET_KEY"`
}

type Server struct {
	ServerPort    string `env:"SERVER_PORT"`
	ServerHost    string `env:"SERVER_HOST"`
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

type Database struct {
	DatabaseDSN string `env:"DATABASE_DSN"`
}

func GetConfig() (Config, error) {
	cfg := Config{
		Server{
			ServerPort:    "8080",
			ServerHost:    "localhost",
			ServerAddress: "localhost:8080",
			BaseURL:       "http://localhost:8080",
		},
		Database{DatabaseDSN: ""},
		Auth{SecretKey: ""},
	}

	// берем конфиг из окружения
	err := env.Parse(&cfg)
	if err != nil {
		return Config{}, err
	}

	// читаем флаги, если есть - перезаписываем конфиг
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "base url")
	flag.StringVar(&cfg.SecretKey, "s", cfg.BaseURL, "secret key")
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "host to listen on")
	flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "database connection string")
	flag.Parse()

	return cfg, nil
}
