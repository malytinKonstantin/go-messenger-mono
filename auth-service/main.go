package main

import (
	"context"
	"log"
	"net"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gofiber/fiber/v2"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/auth-service"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedAuthServiceServer
	producer *kafka.Producer
}

func (s *server) Authenticate(ctx context.Context, in *pb.AuthRequest) (*pb.AuthResponse, error) {
	// Реализация аутентификации
	return &pb.AuthResponse{Token: "sample_token"}, nil
}

func main() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer p.Close()

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, &server{producer: p})

	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Auth Service is healthy")
	})

	go func() {
		log.Fatal(app.Listen(":3002"))
	}()

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
