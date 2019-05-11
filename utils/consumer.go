package utils

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func Consume(topics []string, group string, callback func(message *kafka.Message)) error {

	fmt.Println("Broker: ", broker)
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  broker,
		"group.id":           group,
		"session.timeout.ms": 6000,
		"auto.offset.reset":  "earliest"})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		return err
	}

	fmt.Printf("Created Consumer %v\n", c)

	err = c.SubscribeTopics(topics, nil)

	run := true

	for run == true {
		ev := c.Poll(100)
		if ev == nil {
			continue
		}

		switch e := ev.(type) {
		case *kafka.Message:
			fmt.Printf("Event received on %s: %s\n",
				e.TopicPartition, string(e.Value))
			callback(e)
		case kafka.Error:
			// Errors should generally be considered as informational, the client will try to automatically recover
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
		default:
			// fmt.Printf("Ignored %v\n", e)
		}
	}

	fmt.Printf("Closing consumer\n")
	c.Close()
	return nil
}
