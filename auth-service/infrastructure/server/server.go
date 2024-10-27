package server

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	grpcHandlers "github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/delivery/grpc"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/repository"
	authUsecase "github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/usecase/auth"
	credentialsUsecase "github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/usecase/credentials"
	oauthUsecase "github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/usecase/oauth"
	passwordUsecase "github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/usecase/password"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/auth_service/v1"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func SetupGRPCServer(db *gorm.DB, producer *kafka.Producer) (*grpc.Server, error) {
	server := grpc.NewServer()

	// Инициализация репозиториев
	userRepo := repository.NewUserCredentialsRepository(db)
	oauthRepo := repository.NewOauthAccountRepository(db)
	tokenRepo := repository.NewResetPasswordTokenRepository(db)

	// Инициализация usecase
	jwtSecret := viper.GetString("JWT_SECRET")
	registerUC := authUsecase.NewRegisterUserUsecase(userRepo)
	authenticateUC := authUsecase.NewAuthenticateUserUsecase(userRepo, jwtSecret)
	verifyEmailUC := credentialsUsecase.NewVerifyEmailUsecase(userRepo)
	oauthAuthenticateUC := oauthUsecase.NewOAuthAuthenticateUsecase(userRepo, oauthRepo)
	resetPasswordUC := passwordUsecase.NewResetPasswordUsecase(userRepo, tokenRepo)
	changePasswordUC := passwordUsecase.NewChangePasswordUsecase(userRepo, tokenRepo)

	// Инициализация хендлеров
	authHandler := grpcHandlers.NewAuthHandler(registerUC, authenticateUC, verifyEmailUC)
	oauthHandler := grpcHandlers.NewOAuthHandler(oauthAuthenticateUC)
	passwordHandler := grpcHandlers.NewPasswordHandler(resetPasswordUC, changePasswordUC)

	// Регистрация сервисов
	pb.RegisterAuthServiceServer(server, authHandler)

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
