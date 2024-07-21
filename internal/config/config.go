package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddr  string
	DatabaseURL string
	LogLevel    string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		ServerAddr:  os.Getenv("SERVER_ADDR"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		LogLevel:    os.Getenv("LOG_LEVEL"),
	}, nil
}
