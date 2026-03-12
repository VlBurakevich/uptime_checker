package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppPort   string
	JWTSecret string

	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}

	Kafka struct {
		Addr         string
		TopicTasks   string
		TopicResults string
	}

	Scheduler struct {
		Limit    int
		Interval time.Duration
	}
}

func Load() *Config {
	cfg := &Config{}

	cfg.AppPort = getEnv("APP_PORT", "8080")
	cfg.JWTSecret = getEnv("JWT_SECRET", "super-secret-secret")

	cfg.DB.Host = getEnv("DB_HOST", "localhost")
	cfg.DB.Port = getEnv("DB_PORT", "5432")
	cfg.DB.User = getEnv("DB_USER", "postgres")
	cfg.DB.Password = getEnv("DB_PASSWORD", "root")
	cfg.DB.Name = getEnv("DB_NAME", "uptime_db")

	cfg.Kafka.Addr = getEnv("KAFKA_ADDR", "localhost:9092")
	cfg.Kafka.TopicTasks = getEnv("TOPIC_TASKS", "site.checks")
	cfg.Kafka.TopicResults = getEnv("TOPIC_RESULTS", "check.results")

	cfg.Scheduler.Limit = getEnvInt("SCHEDULER_LIMIT", 50)
	cfg.Scheduler.Interval = getEnvDuration("SCHEDULER_INTERVAL", 10*time.Second)

	return cfg
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.DB.Host, c.DB.User, c.DB.Password, c.DB.Name, c.DB.Port,
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

func getEnvInt(key string, fallback int) int {
	if val, exists := os.LookupEnv(key); exists {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if val, exists := os.LookupEnv(key); exists {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return fallback
}
