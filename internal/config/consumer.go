package config

// Consumer holds the configuration values for the Kafka consumer.
type Consumer struct {
	Brokers []string `env:"KAFKA_BROKERS" default:"[192.168.178.17:9092]"`
	Topic   string   `env:"KAFKA_STORE_CHOICE_CHOICE_TOPIC" default:"store.choice.choice"`
	GroupID string   `env:"KAFKA_CONSUMER_CHOICE_GROUP_ID" default:"consumer-choice-group"`
}
