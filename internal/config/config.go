package config

type Config struct {
	KafkaAddr        string `env:"KAFKA_ADDR" envDefault:"localhost:9092"`
	TopicSiteCheck   string `env:"KAFKA_TOPIC_SITE_CHECK" envDefault:"site.check"`
	TopicCheckResult string `env:"KAFKA_TOPIC_CHECK_RESULT" envDefault:"chec.result"`
}
