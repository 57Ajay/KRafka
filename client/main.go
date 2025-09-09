package main

import (
	"context"
	"log"
	"time"

	pb "github.com/57ajay/krafka/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	grpcAddr = "localhost:50051"
)

func main() {
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewOrderServiceClient(conn)
	orderId := uuid.New().String()
	order := &pb.Order{
		OrderId:   orderId,
		UserId:    "user-123",
		ItemIds:   []string{"item-abc", "item-def"},
		Timestamp: timestamppb.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	log.Printf("Sending CreateOrder request for Order ID: %s", orderId)

	res, err := client.CreateOrder(ctx, &pb.CreateOrderRequest{Order: order})

	if err != nil {
		log.Fatalf("Could not create order: %v", err)
	}

	log.Printf("✅ Order created successfully! Server responded with Order ID: %s", res.GetOrderId())
}
