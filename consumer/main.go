package main

import (
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:29092",
		"group.id":          "my-consumer-group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	defer c.Close()

	topic := "events"
	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %s", err)
	}

	fmt.Println("Consumer is running... waiting for messages.")

	for {
		msg, err := c.ReadMessage(100 * time.Millisecond)

		if err == nil {
			fmt.Printf("✅ Received message from %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else if kafkaErr, ok := err.(kafka.Error); !ok || kafkaErr.Code() != kafka.ErrTimedOut {
			fmt.Printf("❌ Consumer error: %v (%v)\n", err, msg)
		}
	}
}
