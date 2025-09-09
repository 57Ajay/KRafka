package main

import (
	"context"
	"net"

	"fmt"
	"log"

	pb "github.com/57ajay/krafka/proto"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

const (
	kafkaTopic = "orders"
	grpcAddr   = "localhost:50051"
)

type server struct {
	pb.UnimplementedOrderServiceServer
	producer *kafka.Producer
}

func (s *server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	order := req.GetOrder()
	log.Printf("Received CreateOrder request for order ID: %s", order.OrderId)

	orderBytes, err := proto.Marshal(order)
	if err != nil {
		log.Printf("Failed to marshal order: %v", err)
		return nil, fmt.Errorf("failed to marshal order: %w", err)
	}
	var topic string = kafkaTopic
	err = s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: orderBytes,
		Key:   []byte(order.OrderId),
	}, nil)

	if err != nil {
		log.Printf("Failed to produce message to Kafka: %v", err)
		return nil, fmt.Errorf("failed to produce message: %w", err)
	}

	log.Printf("Order %s successfully published to Kafka topic %s", order.OrderId, kafkaTopic)

	return &pb.CreateOrderResponse{OrderId: order.OrderId}, nil
}

func main() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":  "localhost:29092",
		"acks":               "all",
		"enable.idempotence": "true",
	})
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %s", err)
	}
	defer p.Close()

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					log.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterOrderServiceServer(s, &server{producer: p})
	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
