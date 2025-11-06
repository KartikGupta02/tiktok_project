package config

import (
	"os"
)

type Config struct {
	Port  string
	Host  string
	Debug string
}

func Load() *Config {
	return &Config{
		Port:  getEnv("PORT", "8080"),
		Host:  getEnv("HOST", "localhost"),
		Debug: getEnv("DEBUG", "false"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
