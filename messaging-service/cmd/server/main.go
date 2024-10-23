package main

import (
	"fmt"
	"log"

	"github.com/malytinKonstantin/go-messenger-mono/messaging-service/infrastructure/database"
	"github.com/malytinKonstantin/go-messenger-mono/messaging-service/infrastructure/queue"
	"github.com/malytinKonstantin/go-messenger-mono/messaging-service/infrastructure/server"
	"github.com/spf13/viper"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error starting service: %v", err)
	}
}

func run() error {
	viper.AutomaticEnv()

	session, err := database.ConnectToCassandra()
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}
	defer session.Close()

	producer, err := queue.CreateKafkaProducer()
	if err != nil {
		return fmt.Errorf("error creating Kafka producer: %w", err)
	}
	defer producer.Close()

	grpcServer, err := server.SetupGRPCServer(session, producer)
	if err != nil {
		return fmt.Errorf("error setting up gRPC server: %w", err)
	}

	httpServer := server.SetupHTTPServer()

	go func() {
		if err := server.StartHTTPServer(httpServer); err != nil {
			log.Printf("Error starting HTTP server: %v", err)
		}
	}()

	go func() {
		if err := server.StartGRPCServer(grpcServer); err != nil {
			log.Printf("Error starting gRPC server: %v", err)
		}
	}()

	server.WaitForShutdown(httpServer, grpcServer)

	return nil
}
