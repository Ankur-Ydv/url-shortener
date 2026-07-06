package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBSSLMode string

	RedisHost string
	RedisPort string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return &Config{
		DBHost:    getValue("DB_HOST", "localhost"),
		DBPort:    getValue("DB_PORT", "5432"),
		DBUser:    getValue("DB_USER", "postgres"),
		DBPass:    getValue("DB_PASS", "password"),
		DBSSLMode: getValue("DB_SSL_MODE", "disable"),

		RedisHost: getValue("REDIS_HOST", "localhost"),
		RedisPort: getValue("REDIS_PORT", "6379"),
	}, nil
}

func getValue(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
