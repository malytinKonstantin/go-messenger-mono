package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

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
		log.Fatalf("error starting application: %v", err)
	}
}

func run() error {
	viper.AutomaticEnv()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grpcMux := runtime.NewServeMux()

	if err := setupGRPCMux(ctx, grpcMux); err != nil {
		return fmt.Errorf("error setting up gRPC mux: %w", err)
	}

	app := setupFiberApp(grpcMux)

	return startServer(app)
}

func setupGRPCMux(ctx context.Context, grpcMux *runtime.ServeMux) error {
	var wg sync.WaitGroup
	var mu sync.Mutex

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	services := []struct {
		name     string
		register func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
	}{
		{"auth", handlers.RegisterAuthService},
		{"user", handlers.RegisterUserService},
		{"friendship", handlers.RegisterFriendshipService},
		{"messaging", handlers.RegisterMessagingService},
		{"notification", handlers.RegisterNotificationService},
	}

	for _, service := range services {
		wg.Add(1)
		go func(srv struct {
			name     string
			register func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
		}) {
			defer wg.Done()

			var endpoint string
			if viper.GetString("ENV") == "local" {
				port := viper.GetString(fmt.Sprintf("%s_SERVICE_GRPC_PORT", srv.name))
				endpoint = fmt.Sprintf("localhost:%s", port)
			} else {
				endpoint = fmt.Sprintf("%s-service:%s", srv.name, viper.GetString(fmt.Sprintf("%s_SERVICE_GRPC_PORT", srv.name)))
			}

			maxRetries := 3
			retryInterval := 2 * time.Second

			var err error
			for i := 1; i <= maxRetries; i++ {
				if ctx.Err() != nil {
					return
				}
				err = srv.register(ctx, grpcMux, endpoint, opts)
				if err == nil {
					mu.Lock()
					log.Printf("Successfully connected to service %s at %s", srv.name, endpoint)
					mu.Unlock()
					return
				}
				mu.Lock()
				log.Printf("Failed to connect to service %s at %s: %v. Attempt %d of %d", srv.name, endpoint, err, i, maxRetries)
				mu.Unlock()
				time.Sleep(retryInterval)
			}
			mu.Lock()
			log.Printf("Failed to connect to service %s after %d attempts", srv.name, maxRetries)
			mu.Unlock()
		}(service)
	}

	wg.Wait()

	return nil
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
