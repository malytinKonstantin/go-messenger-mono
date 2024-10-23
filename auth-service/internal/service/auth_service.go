package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.AuthRepository
}

func NewAuthService(repo *repository.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

// Регистрирует нового пользователя
func (s *AuthService) RegisterUser(ctx context.Context, email, password string) (*db.UserCredentials, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.CreateUser(ctx, email, string(hashedPassword))
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Аутентифицирует пользователя
func (s *AuthService) AuthenticateUser(ctx context.Context, email, password string) (*db.UserCredentials, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("incorrect password")
	}

	return user, nil
}

// Подтверждает аккаунт пользователя
func (s *AuthService) VerifyUser(ctx context.Context, userID uuid.UUID) error {
	return s.repo.UpdateUserVerificationStatus(ctx, userID, true)
}

// Создает OAuth аккаунт
func (s *AuthService) CreateOAuthAccount(ctx context.Context, userID uuid.UUID, provider, providerUserID string) (*db.OauthAccounts, error) {
	return s.repo.CreateOAuthAccount(ctx, userID, provider, providerUserID)
}

// Получает OAuth аккаунт
func (s *AuthService) GetOAuthAccount(ctx context.Context, provider, providerUserID string) (*db.OauthAccounts, error) {
	return s.repo.GetOAuthAccount(ctx, provider, providerUserID)
}

// Создает токен для сброса пароля
func (s *AuthService) CreateResetPasswordToken(ctx context.Context, userID uuid.UUID) (*db.ResetPasswordTokens, error) {
	expiresAt := time.Now().Add(24 * time.Hour) // Токен действителен 24 часа
	return s.repo.CreateResetPasswordToken(ctx, userID, expiresAt)
}

// Получает токен для сброса пароля
func (s *AuthService) GetResetPasswordToken(ctx context.Context, token uuid.UUID) (*db.ResetPasswordTokens, error) {
	return s.repo.GetResetPasswordToken(ctx, token)
}

// Удаляет токен для сброса пароля
func (s *AuthService) DeleteResetPasswordToken(ctx context.Context, token uuid.UUID) error {
	return s.repo.DeleteResetPasswordToken(ctx, token)
}

// Обновляет пароль пользователя
func (s *AuthService) UpdateUserPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.UpdateUserPassword(ctx, userID, string(hashedPassword))
}
