package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/malytinKonstantin/go-messenger-mono/api-gateway/internal/handlers"
	"github.com/malytinKonstantin/go-messenger-mono/api-gateway/internal/middleware"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

var redisClient *redis.Client

func main() {
	if err := run(); err != nil {
		log.Fatalf("error starting application: %v", err)
	}
}

func run() error {
	viper.AutomaticEnv()

	redisClient = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("REDIS_HOST"),
		Password: viper.GetString("REDIS_PASSWORD"),
		DB:       0,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return fmt.Errorf("Ошибка подключения к Redis: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grpcMux := runtime.NewServeMux()

	if err := setupGRPCMux(ctx, grpcMux); err != nil {
		return fmt.Errorf("error setting up gRPC mux: %w", err)
	}

	app := setupFiberApp(grpcMux)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := startServer(app); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Shutdown(); err != nil {
		return fmt.Errorf("Server shutdown failed: %w", err)
	}

	log.Println("Server exited gracefully")
	return nil
}

func setupGRPCMux(ctx context.Context, grpcMux *runtime.ServeMux) error {
	var wg sync.WaitGroup
	var mu sync.Mutex

	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(100 * time.Millisecond)),
		grpc_retry.WithCodes(codes.Unavailable, codes.DeadlineExceeded),
		grpc_retry.WithMax(3),
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)),
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
	app.Use(middleware.IdempotencyMiddleware(redisClient))
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Second,
	}))
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
