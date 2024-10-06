package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/malytinKonstantin/go-messenger-mono/api-gateway/internal/handlers"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	grpcMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Регистрируем сервисы
	if err := handlers.RegisterAuthService(ctx, grpcMux, fmt.Sprintf("auth-service:%s", viper.GetString("AUTH_SERVICE_GRPC_PORT")), opts); err != nil {
		log.Fatalf("Failed to register auth service: %v", err)
	}
	if err := handlers.RegisterFriendshipService(ctx, grpcMux, fmt.Sprintf("friendship-service:%s", viper.GetString("FRIENDSHIP_SERVICE_GRPC_PORT")), opts); err != nil {
		log.Fatalf("Failed to register friendship service: %v", err)
	}
	if err := handlers.RegisterMessagingService(ctx, grpcMux, fmt.Sprintf("messaging-service:%s", viper.GetString("MESSAGING_SERVICE_GRPC_PORT")), opts); err != nil {
		log.Fatalf("Failed to register messaging service: %v", err)
	}
	if err := handlers.RegisterNotificationService(ctx, grpcMux, fmt.Sprintf("notification-service:%s", viper.GetString("NOTIFICATION_SERVICE_GRPC_PORT")), opts); err != nil {
		log.Fatalf("Failed to register notification service: %v", err)
	}
	if err := handlers.RegisterUserService(ctx, grpcMux, fmt.Sprintf("user-service:%s", viper.GetString("USER_SERVICE_GRPC_PORT")), opts); err != nil {
		log.Fatalf("Failed to register user service: %v", err)
	}

	app := fiber.New()
	app.All("/*", adaptor.HTTPHandler(grpcMux))

	port := viper.GetString("HTTP_PORT")
	log.Printf("Starting API Gateway on :%s", port)
	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
