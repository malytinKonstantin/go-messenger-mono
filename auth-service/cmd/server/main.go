package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/handlers"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/auth_service/v1"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error starting service: %v", err)
	}
}

func run() error {
	viper.AutomaticEnv()

	db, err := connectToDatabase()
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}
	defer db.Close()

	producer, err := createKafkaProducer()
	if err != nil {
		return fmt.Errorf("error creating Kafka producer: %w", err)
	}
	defer producer.Close()

	grpcServer, err := setupGRPCServer(db, producer)
	if err != nil {
		return fmt.Errorf("error setting up gRPC server: %w", err)
	}

	httpServer := setupHTTPServer()

	go func() {
		if err := httpServer.Listen(fmt.Sprintf(":%s", viper.GetString("HTTP_PORT"))); err != nil {
			log.Printf("Error starting HTTP server: %v", err)
		}
	}()

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", viper.GetString("GRPC_PORT")))
		if err != nil {
			log.Printf("Error listening on gRPC port: %v", err)
			return
		}
		log.Printf("gRPC server listening at %v", lis.Addr())
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("Error starting gRPC server: %v", err)
		}
	}()

	waitForShutdown(httpServer, grpcServer)

	return nil
}

func connectToDatabase() (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		viper.GetString("DATABASE_USER"),
		viper.GetString("DATABASE_PASSWORD"),
		viper.GetString("DATABASE_HOST"),
		viper.GetString("DATABASE_PORT"),
		viper.GetString("DATABASE_NAME"))

	return pgxpool.New(context.Background(), connStr)
}

func createKafkaProducer() (*kafka.Producer, error) {
	return kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": viper.GetString("KAFKA_BOOTSTRAP_SERVERS")})
}

func setupGRPCServer(db *pgxpool.Pool, producer *kafka.Producer) (*grpc.Server, error) {
	server := grpc.NewServer()
	authHandler := handlers.NewAuthHandler(producer, db)
	pb.RegisterAuthServiceServer(server, authHandler)
	reflection.Register(server)
	return server, nil
}

func setupHTTPServer() *fiber.App {
	app := fiber.New()
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Auth Service is healthy")
	})
	return app
}

func waitForShutdown(httpServer *fiber.App, grpcServer *grpc.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")
	if err := httpServer.Shutdown(); err != nil {
		log.Printf("Error shutting down HTTP server: %v", err)
	}
	grpcServer.GracefulStop()
	log.Println("Servers shut down")
}
