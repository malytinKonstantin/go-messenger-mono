package main

import (
	"context"
	"log"
	"net"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gofiber/fiber/v2"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/messaging-service"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMessagingServiceServer
	producer *kafka.Producer
}

func (s *server) SendMessage(ctx context.Context, in *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	// Реализация отправки сообщения
	return &pb.SendMessageResponse{Success: true}, nil
}

func main() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer p.Close()

	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMessagingServiceServer(s, &server{producer: p})

	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Messaging Service is healthy")
	})

	go func() {
		log.Fatal(app.Listen(":3004"))
	}()

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
