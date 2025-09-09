package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type EventsData struct {
	UserID    string    `json:"user_id"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":  "localhost:29092",
		"acks":               "all",
		"enable.idempotence": "true",
	})
	if err != nil {
		log.Fatalf("Failed to create producer.\nERROR: %s", err.Error())
	}
	defer p.Close()

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
				} else {
					fmt.Printf("Delivered message with key '%s' to topic %s [%d] at offset %v\n",
						string(ev.Key), *ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}
			}
		}
	}()

	topic := "events"
	users := []string{"user-X", "user-Y", "user-Z"}

	for _ = range 100000 {
		for i := range 10 {
			key := users[i%len(users)]
			event := EventsData{
				UserID:    key,
				Message:   "Event number " + strconv.Itoa(i) + "-" + key,
				Timestamp: time.Now(),
			}

			eventBytes, err := json.Marshal(event)
			if err != nil {
				log.Fatalf("Failed to marshal event: %s", err)
			}

			err = p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Key:            []byte(key),
				Value:          eventBytes,
			}, nil)

			if err != nil {
				log.Fatalf("Failed to produce message: %s", err)
			}
		}
	}
	p.Flush(15 * 1000)

	fmt.Println("Message Sent")
}
