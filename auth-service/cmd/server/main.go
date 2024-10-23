package main

import (
	"fmt"
	"log"
	"net"

	"github.com/malytinKonstantin/go-messenger-mono/auth-service/infrastructure/database"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/infrastructure/queue"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/infrastructure/server"
	"github.com/spf13/viper"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error starting service: %v", err)
	}
}

func run() error {
	viper.AutomaticEnv()

	db, err := database.ConnectToDatabase()
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}
	defer db.Close()

	producer, err := queue.CreateKafkaProducer()
	if err != nil {
		return fmt.Errorf("error creating Kafka producer: %w", err)
	}
	defer producer.Close()

	grpcServer, err := server.SetupGRPCServer(db, producer)
	if err != nil {
		return fmt.Errorf("error setting up gRPC server: %w", err)
	}

	httpServer := server.SetupHTTPServer()

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

	server.WaitForShutdown(httpServer, grpcServer)

	return nil
}
