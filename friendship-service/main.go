package main

import (
	"context"
	"log"
	"net"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gofiber/fiber/v2"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/friendship-service"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedFriendshipServiceServer
	producer *kafka.Producer
}

func (s *server) SendFriendRequest(ctx context.Context, in *pb.SendFriendRequestRequest) (*pb.SendFriendRequestResponse, error) {
	// Реализация отправки запроса на дружбу
	log.Printf("Received friend request from %s to %s", in.SenderId, in.RecipientId)
	// Здесь должна быть логика обработки запроса на дружбу
	return &pb.SendFriendRequestResponse{Success: true}, nil
}

func (s *server) AcceptFriendRequest(ctx context.Context, in *pb.AcceptFriendRequestRequest) (*pb.AcceptFriendRequestResponse, error) {
	// Реализация принятия запроса на дружбу
	log.Printf("Accepting friend request %s by user %s", in.RequestId, in.AccepterId)
	// Здесь должна быть логика обработки принятия запроса на дружбу
	return &pb.AcceptFriendRequestResponse{Success: true}, nil
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
	pb.RegisterFriendshipServiceServer(s, &server{producer: p})

	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Friendship Service is healthy")
	})

	go func() {
		log.Fatal(app.Listen(":3002"))
	}()

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
