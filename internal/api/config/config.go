package config

import (
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	AppPort   string `env:"APP_PORT" envDefault:"8080"`
	JWTSecret string `env:"JWT_SECRET" envDefault:"super-secret-secret"`

	DB struct {
		Host     string `env:"HOST" envDefault:"localhost"`
		Port     string `env:"PORT" envDefault:"5432"`
		User     string `env:"USER" envDefault:"postgres"`
		Password string `env:"PASSWORD" envDefault:"root"`
		Name     string `env:"NAME" envDefault:"uptime_db"`
	} `envPrefix:"DB_"`

	Kafka struct {
		KafkaBroker  string `env:"BROKER" envDefault:"localhost:9092"`
		TopicTasks   string `env:"TOPIC_TASKS" envDefault:"site.checks"`
		TopicResults string `env:"TOPIC_RESULTS" envDefault:"check.results"`
	} `envPrefix:"KAFKA_"`

	Scheduler struct {
		Limit    int           `env:"LIMIT" envDefault:"50"`
		Interval time.Duration `env:"INTERVAL" envDefault:"10s"`
	} `envPrefix:"SCHEDULER_"`
}

func Load() *Config {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	return cfg
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.DB.Host, c.DB.User, c.DB.Password, c.DB.Name, c.DB.Port,
	)
}
