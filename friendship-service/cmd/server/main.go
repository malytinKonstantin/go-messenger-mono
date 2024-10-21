package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/malytinKonstantin/go-messenger-mono/friendship-service/internal/handlers"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/friendship_service/v1"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
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
	if err := loadEnv(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
	viper.AutomaticEnv()

	driver, err := connectToNeo4j()
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}
	defer driver.Close()

	producer, err := createKafkaProducer()
	if err != nil {
		return fmt.Errorf("error creating Kafka producer: %w", err)
	}
	defer producer.Close()

	grpcServer, err := setupGRPCServer(driver, producer)
	if err != nil {
		return fmt.Errorf("error setting up gRPC server: %w", err)
	}

	httpServer := setupHTTPServer()

	go func() {
		if err := startHTTPServer(httpServer); err != nil {
			log.Printf("Error starting HTTP server: %v", err)
		}
	}()

	go func() {
		if err := startGRPCServer(grpcServer); err != nil {
			log.Printf("Error starting gRPC server: %v", err)
		}
	}()

	waitForShutdown(httpServer, grpcServer)

	return nil
}

func loadEnv() error {
	env := os.Getenv("ENV")
	if env != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println(".env file not found, continuing without it")
		}
	}
	return nil
}

func connectToNeo4j() (neo4j.Driver, error) {
	driver, err := neo4j.NewDriver(
		fmt.Sprintf("bolt://%s:%s", viper.GetString("DATABASE_HOST"), viper.GetString("DATABASE_PORT")),
		neo4j.BasicAuth(
			viper.GetString("NEO4J_USER"),
			viper.GetString("NEO4J_PASSWORD"),
			"",
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Neo4j driver: %v", err)
	}
	return driver, nil
}

func createKafkaProducer() (*kafka.Producer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": viper.GetString("KAFKA_BOOTSTRAP_SERVERS"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %v", err)
	}
	return producer, nil
}

func setupGRPCServer(driver neo4j.Driver, producer *kafka.Producer) (*grpc.Server, error) {
	grpcServer := grpc.NewServer()
	friendshipHandler := handlers.NewFriendshipHandler(producer, driver)
	pb.RegisterFriendshipServiceServer(grpcServer, friendshipHandler)
	reflection.Register(grpcServer)
	return grpcServer, nil
}

func setupHTTPServer() *fiber.App {
	app := fiber.New()
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Friendship Service is healthy")
	})
	return app
}

func startHTTPServer(app *fiber.App) error {
	if err := app.Listen(fmt.Sprintf(":%s", viper.GetString("HTTP_PORT"))); err != nil {
		return fmt.Errorf("error starting HTTP server: %v", err)
	}
	return nil
}

func startGRPCServer(grpcServer *grpc.Server) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", viper.GetString("GRPC_PORT")))
	if err != nil {
		return fmt.Errorf("failed to start listener: %v", err)
	}
	log.Printf("gRPC server listening on %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("error starting gRPC server: %v", err)
	}
	return nil
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
	log.Println("Servers successfully shut down")
}
