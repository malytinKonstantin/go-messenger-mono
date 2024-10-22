package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/friendship-service"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedFriendshipServiceServer
	producer *kafka.Producer
}

func (s *server) SendFriendRequest(ctx context.Context, in *pb.SendFriendRequestRequest) (*pb.SendFriendRequestResponse, error) {
	// Реализация отправки запроса на дружбу
	// log.Printf("Received friend request from %s to %s", in.SenderId, in.RecipientId)
	// Здесь должна быть логика обработки запроса на дружбу
	return &pb.SendFriendRequestResponse{Success: true}, nil
}

func (s *server) AcceptFriendRequest(ctx context.Context, in *pb.AcceptFriendRequestRequest) (*pb.AcceptFriendRequestResponse, error) {
	// Реализация принятия запроса на дружбу
	// log.Printf("Accepting friend request %s by user %s", in.RequestId, in.AccepterId)
	// Здесь должна быть логика обработки принятия запроса на дружбу
	return &pb.AcceptFriendRequestResponse{Success: true}, nil
}

func main() {
	env := os.Getenv("ENV")
	if env != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Файл .env не найден, продолжаем без него")
		}
	}
	viper.AutomaticEnv()

	driver, err := neo4j.NewDriver(
		viper.GetString("DATABASE_HOST"),
		neo4j.BasicAuth(
			viper.GetString("NEO4J_USER"),
			viper.GetString("NEO4J_PASSWORD"),
			"",
		),
	)
	if err != nil {
		log.Fatalf("Failed to create Neo4j driver: %v", err)
	}
	defer driver.Close()

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": viper.GetString("KAFKA_BOOTSTRAP_SERVERS")})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer p.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", viper.GetString("GRPC_PORT")))
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
		log.Fatal(app.Listen(fmt.Sprintf(":%s", viper.GetString("HTTP_PORT"))))
	}()

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
