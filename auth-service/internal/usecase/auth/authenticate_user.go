package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticateUserUsecase interface {
	Execute(ctx context.Context, email, password string) (string, error)
}

type authenticateUserUsecase struct {
	userRepo  repository.UserCredentialsRepository
	jwtSecret string
}

func NewAuthenticateUserUsecase(userRepo repository.UserCredentialsRepository, jwtSecret string) AuthenticateUserUsecase {
	return &authenticateUserUsecase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (uc *authenticateUserUsecase) Execute(ctx context.Context, email, password string) (string, error) {
	// Получаем пользователя по email
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	userID, err := uuid.Parse(user.UserID)
	if err != nil {
		return "", errors.New("invalid user_id")
	}

	// Генерируем JWT токен
	token, err := uc.generateJWTToken(userID)
	if err != nil {
		return "", errors.New("error creating token")
	}

	return token, nil
}

func (uc *authenticateUserUsecase) generateJWTToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
