package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	handlers "github.com/malytinKonstantin/go-messenger-mono/friendship-service/internal/delivery/grpc"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gofiber/fiber/v2"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/friendship_service/v1"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func SetupGRPCServer(producer *kafka.Producer, driver neo4j.DriverWithContext) (*grpc.Server, error) {
	server := grpc.NewServer()

	friendshipHandler := handlers.NewFriendshipHandler(producer, driver)

	pb.RegisterFriendshipServiceServer(server, friendshipHandler)

	listener, err := net.Listen("tcp", ":50053")
	if err != nil {
		return nil, err
	}

	go func() {
		if err := server.Serve(listener); err != nil {
			panic(err)
		}
	}()

	return server, nil
}

func SetupHTTPServer() *fiber.App {
	app := fiber.New()
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Friendship Service is healthy")
	})
	return app
}

func StartHTTPServer(app *fiber.App) error {
	if err := app.Listen(fmt.Sprintf(":%s", viper.GetString("HTTP_PORT"))); err != nil {
		return fmt.Errorf("error starting HTTP server: %v", err)
	}
	return nil
}

func StartGRPCServer(grpcServer *grpc.Server) error {
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

func WaitForShutdown(httpServer *fiber.App, grpcServer *grpc.Server) {
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
