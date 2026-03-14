package config

import (
	"log"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	AppPort string `env:"APP_PORT" envDefault:"8080"`

	Kafka struct {
		Broker       string `env:"BROKER" envDefault:"localhost:9092"`
		TopicTasks   string `env:"TOPIC_TASKS" envDefault:"site.checks"`
		TopicResults string `env:"TOPIC_RESULTS" envDefault:"check.results"`
		GroupId      string `env:"GROUP_ID" envDefault:"api-group"`
	} `envPrefix:"KAFKA_"`

	Watcher struct {
		HTTPTimeout time.Duration `env:"HTTP_TIMEOUT" envDefault:"10s"`
		MaxTaskAge  time.Duration `env:"MAX_TASK_AGE" envDefault:"5m"`

		MinConcurrency int32 `env:"MIN_CONCURRENCY" envDefault:"5"`
		MaxConcurrency int32 `env:"MAX_CONCURRENCY" envDefault:"200"`
	}
}

func Load() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal("failed to parse config", err)
	}
	return cfg
}
