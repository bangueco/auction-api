package config

import "os"

type Config struct {
	DATABASE_URL string
	PORT         string
}

func Load() *Config {
	return &Config{
		DATABASE_URL: os.Getenv("DATABASE_URL"),
		PORT:         os.Getenv("PORT"),
	}
}
