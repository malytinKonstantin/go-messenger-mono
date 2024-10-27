package repository

import (
	"context"
	"time"

	"github.com/malytinKonstantin/go-messenger-mono/auth-service/infrastructure/database/generated"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/infrastructure/database/model"
	"gorm.io/gorm"
)

type ResetPasswordTokenRepository interface {
	Create(ctx context.Context, token *model.ResetPasswordToken) error
	GetByToken(ctx context.Context, token string) (*model.ResetPasswordToken, error)
	Delete(ctx context.Context, token string) error
	DeleteExpiredTokens(ctx context.Context) error
}

type resetPasswordTokenRepository struct {
	db *gorm.DB
}

func NewResetPasswordTokenRepository(db *gorm.DB) ResetPasswordTokenRepository {
	return &resetPasswordTokenRepository{
		db: db,
	}
}

func (r *resetPasswordTokenRepository) Create(ctx context.Context, token *model.ResetPasswordToken) error {
	return generated.ResetPasswordToken.WithContext(ctx).Create(token)
}

func (r *resetPasswordTokenRepository) GetByToken(ctx context.Context, token string) (*model.ResetPasswordToken, error) {
	return generated.ResetPasswordToken.WithContext(ctx).Where(generated.ResetPasswordToken.Token.Eq(token)).First()
}

func (r *resetPasswordTokenRepository) Delete(ctx context.Context, token string) error {
	_, err := generated.ResetPasswordToken.WithContext(ctx).Where(generated.ResetPasswordToken.Token.Eq(token)).Delete()
	return err
}

func (r *resetPasswordTokenRepository) DeleteExpiredTokens(ctx context.Context) error {
	_, err := generated.ResetPasswordToken.WithContext(ctx).
		Where(generated.ResetPasswordToken.ExpiresAt.Lt(time.Now())).
		Delete()
	return err
}
