package config

import (
	"os"
)

type Config struct {
	Port       string
	DBSource   string
	JWTSecret  string
}

func LoadConfig() *Config {
	return &Config{
		Port:      getEnv("PORT", "8080"),
		DBSource: getEnv("DB_SOURCE", "postgres://user:password@db:5432/merchstore?sslmode=disable"),
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
