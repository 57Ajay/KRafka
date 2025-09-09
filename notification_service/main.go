package main

import (
	"log"
	"time"

	pb "github.com/57ajay/krafka/proto"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"google.golang.org/protobuf/proto"
)

const (
	kafkaTopic = "orders"
)

func main() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:29092",
		"group.id":          "notification-group",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	defer c.Close()

	err = c.SubscribeTopics([]string{kafkaTopic}, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %s", err)
	}

	log.Println("Notification service is running... waiting for orders.")

	for {
		msg, err := c.ReadMessage(100 * time.Millisecond)
		if err == nil {
			order := &pb.Order{}
			if err := proto.Unmarshal(msg.Value, order); err != nil {
				log.Printf("Failed to unmarshal order: %v", err)
				continue
			}

			log.Printf("📧 [Notification] Sending notification for order %s to user %s.",
				order.OrderId, order.UserId)
			time.Sleep(500 * time.Millisecond)

		} else if kafkaErr, ok := err.(kafka.Error); !ok || kafkaErr.Code() != kafka.ErrTimedOut {
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
