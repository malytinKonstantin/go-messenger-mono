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
	handlers "github.com/malytinKonstantin/go-messenger-mono/notification-service/internal/delivery/grpc"
	"github.com/malytinKonstantin/go-messenger-mono/notification-service/internal/repositories"
	notification_uc "github.com/malytinKonstantin/go-messenger-mono/notification-service/internal/usecase/notification"
	preferences_uc "github.com/malytinKonstantin/go-messenger-mono/notification-service/internal/usecase/preferences"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/notification_service/v1"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SetupGRPCServer(session *gocql.Session, producer *kafka.Producer) (*grpc.Server, error) {
	// Инициализация репозиториев
	notificationRepo := repositories.NewNotificationRepository(session)
	preferencesRepo := repositories.NewNotificationPreferencesRepository(session)

	// Инициализация usecase
	sendNotificationUsecase := notification_uc.NewSendNotificationUsecase(notificationRepo)
	getNotificationsUsecase := notification_uc.NewGetNotificationsUsecase(notificationRepo)
	markAsReadUsecase := notification_uc.NewMarkNotificationAsReadUsecase(notificationRepo)
	updatePreferencesUsecase := preferences_uc.NewUpdatePreferencesUsecase(preferencesRepo)
	getPreferencesUsecase := preferences_uc.NewGetPreferencesUsecase(preferencesRepo)

	// Создание gRPC сервера
	grpcServer := grpc.NewServer()

	// Инициализация обработчика сервиса уведомлений
	notificationHandler := handlers.NewNotificationServiceServer(
		sendNotificationUsecase,
		getNotificationsUsecase,
		markAsReadUsecase,
		updatePreferencesUsecase,
		getPreferencesUsecase,
	)

	// Регистрация сервиса
	pb.RegisterNotificationServiceServer(grpcServer, notificationHandler)

	// Включение отражения сервера (для инструментов типа grpcurl)
	reflection.Register(grpcServer)

	return grpcServer, nil
}

func SetupHTTPServer() *fiber.App {
	app := fiber.New()
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Notification Service is healthy")
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
