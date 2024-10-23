package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	db "github.com/malytinKonstantin/go-messenger-mono/auth-service/infrastructure/database/generated"
)

type AuthRepository struct {
	q *db.Queries
}

func NewAuthRepository(dbConn *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		q: db.New(stdlib.OpenDBFromPool(dbConn)),
	}
}

// CreateUser создает нового пользователя
func (r *AuthRepository) CreateUser(ctx context.Context, email, passwordHash string) (*db.UserCredentials, error) {
	params := db.CreateUserParams{
		Email:        email,
		PasswordHash: passwordHash,
	}
	user, err := r.q.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail получает пользователя по email
func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*db.UserCredentials, error) {
	user, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUserVerificationStatus обновляет статус верификации пользователя
func (r *AuthRepository) UpdateUserVerificationStatus(ctx context.Context, userID uuid.UUID, isVerified bool) error {
	params := db.UpdateUserVerificationStatusParams{
		UserID:     userID,
		IsVerified: isVerified,
	}
	return r.q.UpdateUserVerificationStatus(ctx, params)
}

// CreateOAuthAccount создает новую запись OAuth аккаунта
func (r *AuthRepository) CreateOAuthAccount(ctx context.Context, userID uuid.UUID, provider, providerUserID string) (*db.OauthAccounts, error) {
	params := db.CreateOAuthAccountParams{
		UserID:         userID,
		Provider:       provider,
		ProviderUserID: providerUserID,
	}
	oauthAccount, err := r.q.CreateOAuthAccount(ctx, params)
	if err != nil {
		return nil, err
	}
	return &oauthAccount, nil
}

// GetOAuthAccount получает OAuth аккаунт по провайдеру и ID пользователя провайдера
func (r *AuthRepository) GetOAuthAccount(ctx context.Context, provider, providerUserID string) (*db.OauthAccounts, error) {
	params := db.GetOAuthAccountParams{
		Provider:       provider,
		ProviderUserID: providerUserID,
	}
	oauthAccount, err := r.q.GetOAuthAccount(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &oauthAccount, nil
}

// CreateResetPasswordToken создает новый токен для сброса пароля
func (r *AuthRepository) CreateResetPasswordToken(ctx context.Context, userID uuid.UUID, expiresAt time.Time) (*db.ResetPasswordTokens, error) {
	expiresAtPtr := &expiresAt
	params := db.CreateResetPasswordTokenParams{
		UserID:    userID,
		ExpiresAt: &expiresAtPtr,
	}
	token, err := r.q.CreateResetPasswordToken(ctx, params)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// GetResetPasswordToken получает токен для сброса пароля
func (r *AuthRepository) GetResetPasswordToken(ctx context.Context, token uuid.UUID) (*db.ResetPasswordTokens, error) {
	resetToken, err := r.q.GetResetPasswordToken(ctx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &resetToken, nil
}

// DeleteResetPasswordToken удаляет токен для сброса пароля
func (r *AuthRepository) DeleteResetPasswordToken(ctx context.Context, token uuid.UUID) error {
	return r.q.DeleteResetPasswordToken(ctx, token)
}

// UpdateUserPassword обновляет пароль пользователя
func (r *AuthRepository) UpdateUserPassword(ctx context.Context, userID uuid.UUID, newPasswordHash string) error {
	params := db.UpdateUserPasswordParams{
		UserID:       userID,
		PasswordHash: newPasswordHash,
	}
	return r.q.UpdateUserPassword(ctx, params)
}
