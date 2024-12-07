package queue

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/spf13/viper"
)

func CreateKafkaProducer() (*kafka.Producer, error) {
	return kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": viper.GetString("KAFKA_BOOTSTRAP_SERVERS")})
}
