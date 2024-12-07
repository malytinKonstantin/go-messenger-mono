package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/user_service/v1"
	handlers "github.com/malytinKonstantin/go-messenger-mono/user-service/internal/delivery/grpc"
	"github.com/malytinKonstantin/go-messenger-mono/user-service/internal/repositories"
	"github.com/malytinKonstantin/go-messenger-mono/user-service/internal/usecases/user"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SetupGRPCServer(session *gocql.Session) (*grpc.Server, error) {
	// Инициализация репозитория
	userRepo := repositories.NewUserRepository(session)

	// Инициализация usecase
	getUserUC := user.NewGetUserUsecase(userRepo)
	createUserUC := user.NewCreateUserProfileUsecase(userRepo)
	updateUserUC := user.NewUpdateUserProfileUsecase(userRepo)
	searchUsersUC := user.NewSearchUsersUsecase(userRepo)

	// Инициализация gRPC хендлера
	userHandler := handlers.NewUserHandler(getUserUC, createUserUC, updateUserUC, searchUsersUC)

	// Создание gRPC сервера и регистрация сервисов
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userHandler)
	reflection.Register(grpcServer)

	return grpcServer, nil
}

func StartGRPCServer(grpcServer *grpc.Server) error {
	port := viper.GetString("GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}

	fmt.Printf("gRPC server started on port %s\n", port)
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("error running gRPC server: %v", err)
	}

	return nil
}

func SetupHTTPServer() *fiber.App {
	app := fiber.New()
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("User Service is healthy")
	})
	return app
}

func StartHTTPServer(app *fiber.App) error {
	port := viper.GetString("HTTP_PORT")
	if port == "" {
		port = "8080"
	}
	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		return fmt.Errorf("error starting HTTP server: %v", err)
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
