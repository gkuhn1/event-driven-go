package utils

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	kafkaProducer *kafka.Producer
}

func NewProducer() *Producer {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": getBroker()})
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		return nil
	}

	fmt.Printf("Created Producer %v\n", producer)
	return &Producer{
		kafkaProducer: producer,
	}
}

func (p *Producer) ProduceMessage(value string, topic string) error {

	err := p.kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(value),
	}, nil)

	fmt.Printf("Produced to %s: %s\n", topic, value)

	return err
}

func (p *Producer) Close() {
	p.kafkaProducer.Close()
}

func getBroker() string {
	return fmt.Sprintf("%s:9092", os.Getenv("BROKER_HOST"))
}
