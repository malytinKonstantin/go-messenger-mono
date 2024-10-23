package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/user_service/v1"
	handlers "github.com/malytinKonstantin/go-messenger-mono/user-service/internal/delivery/grpc"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SetupGRPCServer(session *gocql.Session, producer *kafka.Producer) (*grpc.Server, error) {
	grpcServer := grpc.NewServer()
	userHandler := handlers.NewUserHandler(session, producer)
	pb.RegisterUserServiceServer(grpcServer, userHandler)
	reflection.Register(grpcServer)
	return grpcServer, nil
}

func SetupHTTPServer() *fiber.App {
	app := fiber.New()
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("User service is running")
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
		return fmt.Errorf("failed to listen: %v", err)
	}
	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("error starting gRPC server: %v", err)
	}
	return nil
}

func WaitForShutdown(httpServer *fiber.App, grpcServer *grpc.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down servers...")

	if err := httpServer.Shutdown(); err != nil {
		log.Printf("error shutting down HTTP server: %v", err)
	}

	grpcServer.GracefulStop()
	log.Println("servers successfully shut down")
}
