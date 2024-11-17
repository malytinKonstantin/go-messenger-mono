package repository

import (
	"context"
	"errors"
	"log"

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
	if db == nil {
		log.Fatal("Database connection is nil")
	}
	return &userCredentialsRepository{
		db: db,
	}
}

func (r *userCredentialsRepository) Create(ctx context.Context, user *model.UserCredential) error {
	userModel := generated.Use(r.db).UserCredential

	generatedUser := &model.UserCredential{
		UserID:       user.UserID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		IsVerified:   user.IsVerified,
	}

	err := userModel.WithContext(ctx).Create(generatedUser)
	if err != nil {
		return err
	}

	return nil
}

func (r *userCredentialsRepository) GetByID(ctx context.Context, userID uuid.UUID) (*model.UserCredential, error) {
	return generated.UserCredential.WithContext(ctx).Where(generated.UserCredential.UserID.Eq(userID.String())).First()
}

func (r *userCredentialsRepository) GetByEmail(ctx context.Context, email string) (*model.UserCredential, error) {
	userModel := generated.Use(r.db).UserCredential

	generatedUser, err := userModel.WithContext(ctx).Where(userModel.Email.Eq(email)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	user := &model.UserCredential{
		UserID:       generatedUser.UserID,
		Email:        generatedUser.Email,
		PasswordHash: generatedUser.PasswordHash,
		IsVerified:   generatedUser.IsVerified,
	}

	return user, nil
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
