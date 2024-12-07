package main

import (
	"fmt"
	"log"

	"github.com/malytinKonstantin/go-messenger-mono/user-service/infrastructure/database"
	"github.com/malytinKonstantin/go-messenger-mono/user-service/infrastructure/queue"
	"github.com/malytinKonstantin/go-messenger-mono/user-service/infrastructure/server"
	"github.com/spf13/viper"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error starting service: %v", err)
	}
}

func run() error {
	viper.AutomaticEnv()

	err := database.ConnectToScylla()
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}
	defer database.ScyllaSession.Close()

	if err := database.RunMigrations("./infrastructure/database/migrations"); err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}

	producer, err := queue.CreateKafkaProducer()
	if err != nil {
		return fmt.Errorf("error creating Kafka producer: %w", err)
	}
	defer producer.Close()

	grpcServer, err := server.SetupGRPCServer(database.ScyllaSession)
	if err != nil {
		return fmt.Errorf("error setting up gRPC server: %w", err)
	}

	// Настройка и запуск HTTP сервера
	httpServer := server.SetupHTTPServer()
	go func() {
		if err := server.StartHTTPServer(httpServer); err != nil {
			log.Fatalf("Error starting HTTP server: %v", err)
		}
	}()

	// Запуск gRPC сервера в отдельной горутине
	go func() {
		if err := server.StartGRPCServer(grpcServer); err != nil {
			log.Fatalf("Error starting gRPC server: %v", err)
		}
	}()

	// Ожидание сигналов завершения
	server.WaitForShutdown(httpServer, grpcServer)

	return nil
}
