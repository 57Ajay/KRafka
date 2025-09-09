package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type EventData struct {
	UserID    string    `json:"user_id"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

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
			var event EventData
			err := json.Unmarshal(msg.Value, &event)
			if err != nil {
				fmt.Printf("Error deserializing message value: %s\n", err)
				continue
			}
			fmt.Printf("✅ Received Event: UserID=%s, Message='%s', Time=%s (from partition %d)\n",
				event.UserID, event.Message, event.Timestamp.Format(time.RFC3339), msg.TopicPartition.Partition)

		} else if kafkaErr, ok := err.(kafka.Error); !ok || kafkaErr.Code() != kafka.ErrTimedOut {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
