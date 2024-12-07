package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/infrastructure/database/generated"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/infrastructure/database/model"
	"gorm.io/gorm"
)

type UserCredentialsRepository interface {
	Create(ctx context.Context, user *model.UserCredential) error
	GetByID(ctx context.Context, userID uuid.UUID) (*model.UserCredential, error)
	GetByEmail(ctx context.Context, email string) (*model.UserCredential, error)
	Update(ctx context.Context, user *model.UserCredential) error
	Delete(ctx context.Context, userID uuid.UUID) error
	UpdateVerificationStatus(ctx context.Context, userID uuid.UUID, isVerified bool) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, newPasswordHash string) error
}

type userCredentialsRepository struct {
	db *gorm.DB
}

func NewUserCredentialsRepository(db *gorm.DB) UserCredentialsRepository {
	return &userCredentialsRepository{
		db: db,
	}
}

func (r *userCredentialsRepository) Create(ctx context.Context, user *model.UserCredential) error {
	return generated.UserCredential.WithContext(ctx).Create(user)
}

func (r *userCredentialsRepository) GetByID(ctx context.Context, userID uuid.UUID) (*model.UserCredential, error) {
	return generated.UserCredential.WithContext(ctx).Where(generated.UserCredential.UserID.Eq(userID.String())).First()
}

func (r *userCredentialsRepository) GetByEmail(ctx context.Context, email string) (*model.UserCredential, error) {
	return generated.UserCredential.WithContext(ctx).Where(generated.UserCredential.Email.Eq(email)).First()
}

func (r *userCredentialsRepository) Update(ctx context.Context, user *model.UserCredential) error {
	_, err := generated.UserCredential.WithContext(ctx).
		Where(generated.UserCredential.ID.Eq(user.ID)).
		Updates(user)
	return err
}

func (r *userCredentialsRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	_, err := generated.UserCredential.WithContext(ctx).Where(generated.UserCredential.UserID.Eq(userID.String())).Delete()
	return err
}

func (r *userCredentialsRepository) UpdateVerificationStatus(ctx context.Context, userID uuid.UUID, isVerified bool) error {
	return r.db.WithContext(ctx).
		Model(&model.UserCredential{}).
		Where("user_id = ?", userID).
		Update("is_verified", isVerified).Error
}

func (r *userCredentialsRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, newPasswordHash string) error {
	return r.db.WithContext(ctx).
		Model(&model.UserCredential{}).
		Where("user_id = ?", userID).
		Update("password_hash", newPasswordHash).Error
}
