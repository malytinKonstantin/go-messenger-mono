package server

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gofiber/fiber/v2"
	grpcHandlers "github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/delivery/grpc"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/repository"
	authUsecase "github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/usecase/auth"
	oauthUsecase "github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/usecase/oauth"
	passwordUsecase "github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/usecase/password"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/auth_service/v1"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func SetupGRPCServer(db *gorm.DB, producer *kafka.Producer) (*grpc.Server, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	server := grpc.NewServer()

	userRepo := repository.NewUserCredentialsRepository(db)
	if userRepo == nil {
		return nil, fmt.Errorf("failed to create user repository")
	}

	oauthRepo := repository.NewOauthAccountRepository(db)
	if oauthRepo == nil {
		return nil, fmt.Errorf("failed to create oauth repository")
	}

	tokenRepo := repository.NewResetPasswordTokenRepository(db)
	if tokenRepo == nil {
		return nil, fmt.Errorf("failed to create token repository")
	}

	// Инициализация usecase
	jwtSecret := viper.GetString("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is not set")
	}
	registerUC := authUsecase.NewRegisterUserUsecase(userRepo)
	authenticateUC := authUsecase.NewAuthenticateUserUsecase(userRepo, jwtSecret)
	verifyEmailUC := authUsecase.NewVerifyEmailUsecase(userRepo)
	oauthAuthenticateUC := oauthUsecase.NewOAuthAuthenticateUsecase(userRepo, oauthRepo)
	resetPasswordRequestUC := passwordUsecase.NewResetPasswordRequestUsecase(userRepo, tokenRepo)
	changePasswordUC := passwordUsecase.NewChangePasswordUsecase(userRepo, tokenRepo)

	// Инициализация хендлеров
	authHandler := grpcHandlers.NewAuthHandler(registerUC, authenticateUC, verifyEmailUC)
	oauthHandler := grpcHandlers.NewOAuthHandler(oauthAuthenticateUC)
	passwordHandler := grpcHandlers.NewPasswordHandler(resetPasswordRequestUC, changePasswordUC)

	// Создание составного обработчика
	combinedHandler := grpcHandlers.NewCombinedAuthHandler(authHandler, oauthHandler, passwordHandler)

	// Регистрация сервиса с использованием составного обработчика
	pb.RegisterAuthServiceServer(server, combinedHandler)

	return server, nil
}

func SetupHTTPServer() *fiber.App {
	app := fiber.New()
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("auth service is healthy")
	})
	return app
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
	log.Println("servers shut down")
}
