package main

import (
	"fmt"
	"log"

	"github.com/malytinKonstantin/go-messenger-mono/friendship-service/infrastructure/database"
	"github.com/malytinKonstantin/go-messenger-mono/friendship-service/infrastructure/queue"
	"github.com/malytinKonstantin/go-messenger-mono/friendship-service/infrastructure/server"
	"github.com/spf13/viper"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error starting service: %v", err)
	}
}

func run() error {
	viper.AutomaticEnv()

	if err := database.InitGogm(); err != nil {
		return fmt.Errorf("error initializing GOGM: %w", err)
	}
	defer database.Gogm.Close()

	if err := database.RunMigrations(); err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}

	producer, err := queue.CreateKafkaProducer()
	if err != nil {
		return fmt.Errorf("error creating Kafka producer: %w", err)
	}
	defer producer.Close()

	grpcServer, err := server.SetupGRPCServer(driver, producer)
	if err != nil {
		return fmt.Errorf("error setting up gRPC server: %w", err)
	}

	httpServer := server.SetupHTTPServer()

	go func() {
		if err := server.StartHTTPServer(httpServer); err != nil {
			log.Printf("error starting HTTP server: %v", err)
		}
	}()

	go func() {
		if err := server.StartGRPCServer(grpcServer); err != nil {
			log.Printf("error starting gRPC server: %v", err)
		}
	}()

	server.WaitForShutdown(httpServer, grpcServer)

	return nil
}
