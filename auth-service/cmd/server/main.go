package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/auth-service"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedAuthServiceServer
	producer *kafka.Producer
	db       *pgxpool.Pool
}

func (s *server) Authenticate(ctx context.Context, in *pb.AuthRequest) (*pb.AuthResponse, error) {
	// Реализация аутентификации
	return &pb.AuthResponse{Token: "sample_token"}, nil
}

func main() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		viper.GetString("DATABASE_USER"),
		viper.GetString("DATABASE_PASSWORD"),
		viper.GetString("DATABASE_HOST"),
		viper.GetString("DATABASE_PORT"),
		viper.GetString("DATABASE_NAME"))

	db, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": viper.GetString("KAFKA_BOOTSTRAP_SERVERS")})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer p.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", viper.GetString("GRPC_PORT")))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, &server{producer: p, db: db})

	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Auth Service is healthy")
	})

	go func() {
		log.Fatal(app.Listen(fmt.Sprintf(":%s", viper.GetString("HTTP_PORT"))))
	}()

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
