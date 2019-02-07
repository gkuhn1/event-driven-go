package utils

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var producer *kafka.Producer
var broker string

func init() {
	broker = fmt.Sprintf("%s:9092", os.Getenv("BROKER_HOST"))
	var err error
	producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		return
	}

	fmt.Printf("Created Producer %v\n", producer)
}

func ProduceMessage(value string, topic string) error {

	err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(value),
	}, nil)

	fmt.Printf("Produced to %s: %s\n", topic, value)

	return err
}
