package config

import (
	"fmt"
	"os"
)

func GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "root"),
		getEnv("DB_NAME", "uptime_db"),
		getEnv("DB_PORT", "5432"),
	)
}

func GetJWTSecret() string {
	return getEnv("JWT_SECRET", "super-secret-secret")
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
