package config

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	AppPort   string        `env:"APP_PORT" envDefault:"8080"`
	JWTSecret string        `env:"JWT_SECRET" envDefault:"super-secret-secret"`
	TokenTTL  time.Duration `env:"TOKEN_TTL" envDefault:"1h"`

	DB struct {
		Host     string `env:"HOST" envDefault:"localhost"`
		Port     string `env:"PORT" envDefault:"5432"`
		User     string `env:"USER" envDefault:"postgres"`
		Password string `env:"PASSWORD" envDefault:"root"`
		Name     string `env:"NAME" envDefault:"uptime_db"`
		SSLMode  string `env:"SSL_MODE" envDefault:"disable"`
	} `envPrefix:"DB_"`

	Kafka struct {
		Broker       string `env:"BROKER" envDefault:"localhost:9092"`
		TopicTasks   string `env:"TOPIC_TASKS" envDefault:"site.checks.task"`
		TopicResults string `env:"TOPIC_RESULTS" envDefault:"site.check.results"`
		GroupId      string `env:"GROUP_ID" envDefault:"api-group"`
	} `envPrefix:"KAFKA_"`

	Scheduler struct {
		Limit    int           `env:"LIMIT" envDefault:"50"`
		Interval time.Duration `env:"INTERVAL" envDefault:"10s"`
	} `envPrefix:"SCHEDULER_"`
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.DB.Host, c.DB.User, c.DB.Password, c.DB.Name, c.DB.Port, c.DB.SSLMode,
	)
}

func Load() *Config {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		slog.Error("failed to parse config", "error", err)
		panic(err)
	}

	return cfg
}
