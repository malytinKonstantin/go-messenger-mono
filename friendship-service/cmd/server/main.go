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
		log.Fatalf("error starting service: %v", err)
	}
}

func run() error {
	viper.AutomaticEnv()

	if err := database.InitNeo4jDriver(); err != nil {
		return fmt.Errorf("error initializing Neo4j driver: %w", err)
	}
	defer database.CloseNeo4jDriver()

	if err := database.RunMigrations(); err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}

	producer, err := queue.CreateKafkaProducer()
	if err != nil {
		return fmt.Errorf("error creating Kafka producer: %w", err)
	}
	defer producer.Close()

	grpcServer, err := server.SetupGRPCServer(producer, database.Neo4jDriver)
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
