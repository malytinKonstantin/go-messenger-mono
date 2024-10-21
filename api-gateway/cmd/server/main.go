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
	if err := run(); err != nil {
		log.Fatalf("Error starting application: %v", err)
	}
}

func run() error {

	viper.AutomaticEnv()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grpcMux, err := setupGRPCMux(ctx)
	if err != nil {
		return fmt.Errorf("error setting up gRPC multiplexer: %w", err)
	}

	app := setupFiberApp(grpcMux)

	return startServer(app)
}

func setupGRPCMux(ctx context.Context) (*runtime.ServeMux, error) {
	grpcMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	services := []struct {
		name     string
		register func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
	}{
		{"auth", handlers.RegisterAuthService},
		{"friendship", handlers.RegisterFriendshipService},
		// {"messaging", handlers.RegisterMessagingService},
		// {"notification", handlers.RegisterNotificationService},
		// {"user", handlers.RegisterUserService},
	}

	for _, service := range services {
		var endpoint string
		if viper.GetString("ENV") == "local" {
			port := viper.GetString(fmt.Sprintf("%s_SERVICE_GRPC_PORT", service.name))
			endpoint = fmt.Sprintf("localhost:%s", port)
		} else {
			endpoint = fmt.Sprintf("%s-service:%s", service.name, viper.GetString(fmt.Sprintf("%s_SERVICE_GRPC_PORT", service.name)))
		}
		if err := service.register(ctx, grpcMux, endpoint, opts); err != nil {
			return nil, fmt.Errorf("error registering %s service: %w", service.name, err)
		}
		log.Printf("Successfully connected to %s service at %s", service.name, endpoint)
	}

	return grpcMux, nil
}

func setupFiberApp(grpcMux *runtime.ServeMux) *fiber.App {
	app := fiber.New()
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
	app.All("/*", adaptor.HTTPHandler(grpcMux))
	return app
}

func startServer(app *fiber.App) error {
	port := viper.GetString("HTTP_PORT")
	log.Printf("Starting API Gateway on port :%s", port)
	return app.Listen(fmt.Sprintf(":%s", port))
}
