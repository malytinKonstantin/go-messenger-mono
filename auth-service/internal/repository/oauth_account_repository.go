package repository

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/infrastructure/database/generated"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/infrastructure/database/model"
	"gorm.io/gorm"
)

type OauthAccountRepository interface {
	Create(ctx context.Context, account *model.OauthAccount) error
	GetByID(ctx context.Context, userID uuid.UUID) (*model.OauthAccount, error)
	GetByProvider(ctx context.Context, provider string, providerUserID string) (*model.OauthAccount, error)
	Update(ctx context.Context, account *model.OauthAccount) error
	Delete(ctx context.Context, userID uuid.UUID) error
}

type oauthAccountRepository struct {
	db *gorm.DB
}

func NewOauthAccountRepository(db *gorm.DB) OauthAccountRepository {
	if db == nil {
		log.Fatal("Database connection is nil")
	}
	return &oauthAccountRepository{
		db: db,
	}
}

func (r *oauthAccountRepository) Create(ctx context.Context, account *model.OauthAccount) error {
	return generated.OauthAccount.WithContext(ctx).Create(account)
}

func (r *oauthAccountRepository) GetByID(ctx context.Context, userID uuid.UUID) (*model.OauthAccount, error) {
	return generated.OauthAccount.WithContext(ctx).Where(generated.OauthAccount.UserID.Eq(userID.String())).First()
}

func (r *oauthAccountRepository) GetByProvider(ctx context.Context, provider string, providerUserID string) (*model.OauthAccount, error) {
	return generated.OauthAccount.WithContext(ctx).
		Where(
			generated.OauthAccount.Provider.Eq(provider),
			generated.OauthAccount.ProviderUserID.Eq(providerUserID),
		).First()
}

func (r *oauthAccountRepository) Update(ctx context.Context, account *model.OauthAccount) error {
	_, err := generated.OauthAccount.WithContext(ctx).
		Where(generated.OauthAccount.ID.Eq(account.ID)).
		Updates(account)
	return err
}

func (r *oauthAccountRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	_, err := generated.OauthAccount.WithContext(ctx).Where(generated.OauthAccount.UserID.Eq(userID.String())).Delete()
	return err
}
