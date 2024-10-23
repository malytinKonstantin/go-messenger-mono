package server

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	handlers "github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/delivery/grpc"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/auth_service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SetupGRPCServer(db *pgxpool.Pool, producer *kafka.Producer) (*grpc.Server, error) {
	server := grpc.NewServer()
	authHandler := handlers.NewAuthHandler(producer, db)
	pb.RegisterAuthServiceServer(server, authHandler)
	reflection.Register(server)
	return server, nil
}

func SetupHTTPServer() *fiber.App {
	app := fiber.New()
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Auth Service is healthy")
	})
	return app
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
	log.Println("Servers shut down")
}
