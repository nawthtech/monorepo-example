package config

import (
	"os"
	"strings"
)

type Config struct {
	Port              string
	CORSAllowedOrigins []string
	DatabaseURL       string
	JWTSecret         string
	Environment       string
	APIVersion        string
}

func LoadConfig() Config {
	origins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")

	return Config{
		Port:              getEnv("PORT", "8080"),
		CORSAllowedOrigins: origins,
		DatabaseURL:       os.Getenv("DATABASE_URL"),
		JWTSecret:         os.Getenv("JWT_SECRET"),
		Environment:       getEnv("ENVIRONMENT", "production"),
		APIVersion:        getEnv("API_VERSION", "v1"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}