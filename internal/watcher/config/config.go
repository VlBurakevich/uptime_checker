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

		Monitor struct {
			CpuLow    float64       `env:"CPU_LOW" envDefault:"50"`
			CpuHigh   float64       `env:"CPU_HIGH" envDefault:"80"`
			Increment int32         `env:"INCREMENT" envDefault:"5"`
			Throttle  float64       `env:"THROTTLE" envDefault:"0.8"`
			Interval  time.Duration `env:"INTERVAL" envDefault:"10s"`
		} `envPrefix:"MONITOR_"`
	}
}

func Load() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal("failed to parse config", err)
	}
	return cfg
}
