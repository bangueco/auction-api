package config

import "os"

type Config struct {
	DATABASE_URL      string
	DATABASE_NAME     string
	DATABASE_USER     string
	DATABASE_PASSWORD string
	BASE_URL          string
	PORT              string
	TOKEN_SECRET      string
}

func Load() *Config {
	return &Config{
		DATABASE_URL:      os.Getenv("DATABASE_URL"),
		PORT:              os.Getenv("PORT"),
		DATABASE_NAME:     os.Getenv("DATABASE_NAME"),
		DATABASE_USER:     os.Getenv("DATABASE_USER"),
		DATABASE_PASSWORD: os.Getenv("DATABASE_PASSWORD"),
		BASE_URL:          os.Getenv("BASE_URL"),
		TOKEN_SECRET:      os.Getenv("TOKEN_SECRET"),
	}
}
