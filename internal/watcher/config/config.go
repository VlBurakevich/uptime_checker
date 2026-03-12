package config

import (
	"log"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	AppPort          string        `env:"APP_PORT" envDefault:"8080"`
	KafkaBroker      string        `env:"KAFKA_BROKER" envDefault:"localhost:9092"`
	TopicSiteTask    string        `env:"KAFKA_TOPIC_SITE_CHECK" envDefault:"site.check.task"`
	TopicCheckResult string        `env:"KAFKA_TOPIC_CHECK_RESULT" envDefault:"site.check.result"`
	HTTPTimeout      time.Duration `env:"HTTP_TIMEOUT" envDefault:"10s"`
}

func Load() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal("failed to parse config", err)
	}
	return cfg
}
