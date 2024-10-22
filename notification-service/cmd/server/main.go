package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/malytinKonstantin/go-messenger-mono/notification-service/internal/handlers"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/notification_service/v1"
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

	session, err := connectToCassandra()
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}
	defer session.Close()

	producer, err := createKafkaProducer()
	if err != nil {
		return fmt.Errorf("error creating Kafka producer: %w", err)
	}
	defer producer.Close()

	grpcServer, err := setupGRPCServer(session, producer)
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

func connectToCassandra() (*gocql.Session, error) {
	host := viper.GetString("CASSANDRA_HOST")
	portStr := viper.GetString("CASSANDRA_PORT")
	keyspace := viper.GetString("CASSANDRA_KEYSPACE")
	username := viper.GetString("CASSANDRA_USERNAME")
	password := viper.GetString("CASSANDRA_PASSWORD")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid Cassandra port number: %v", err)
	}

	cluster := gocql.NewCluster(host)
	cluster.Port = port
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum

	if username != "" && password != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: username,
			Password: password,
		}
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Cassandra: %v", err)
	}
	return session, nil
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

func setupGRPCServer(session *gocql.Session, producer *kafka.Producer) (*grpc.Server, error) {
	grpcServer := grpc.NewServer()
	notificationHandler := handlers.NewNotificationServiceServer(producer, session)
	pb.RegisterNotificationServiceServer(grpcServer, notificationHandler)
	reflection.Register(grpcServer)
	return grpcServer, nil
}

func setupHTTPServer() *fiber.App {
	app := fiber.New()
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Notification Service is healthy")
	})
	return app
}

func startHTTPServer(app *fiber.App) error {
	if err := app.Listen(fmt.Sprintf(":%s", viper.GetString("HTTP_PORT"))); err != nil {
		return fmt.Errorf("Error starting HTTP server: %v", err)
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
		return fmt.Errorf("Error starting gRPC server: %v", err)
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
