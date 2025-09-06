package config

import (
	"os"
)

type Config struct {
	Port         string
	DatabaseURL  string
	GeminiAPIKey string
}

func Load() (*Config, error) {
	config := &Config{
		Port:         getEnv("PORT", "8080"),
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/company_ai?sslmode=disable"),
		GeminiAPIKey: getEnv("GEMINI_API_KEY", ""),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
