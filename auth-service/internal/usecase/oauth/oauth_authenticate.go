package oauth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/infrastructure/database/model"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/repository"
)

type OAuthAuthenticateUsecase interface {
	Execute(ctx context.Context, provider, providerUserID, email string) (*model.UserCredential, error)
}

type oauthAuthenticateUsecase struct {
	userRepo  repository.UserCredentialsRepository
	oauthRepo repository.OauthAccountRepository
}

func NewOAuthAuthenticateUsecase(userRepo repository.UserCredentialsRepository, oauthRepo repository.OauthAccountRepository) OAuthAuthenticateUsecase {
	return &oauthAuthenticateUsecase{
		userRepo:  userRepo,
		oauthRepo: oauthRepo,
	}
}

func (uc *oauthAuthenticateUsecase) Execute(ctx context.Context, provider, providerUserID, email string) (*model.UserCredential, error) {
	// Проверяем, существует ли OAuth аккаунт
	oauthAccount, err := uc.oauthRepo.GetByProvider(ctx, provider, providerUserID)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	var user *model.UserCredential

	if oauthAccount != nil {
		// Получаем пользователя по UserID
		userID, err := uuid.Parse(oauthAccount.UserID)
		if err != nil {
			return nil, err
		}

		user, err = uc.userRepo.GetByID(ctx, userID)
		if err != nil {
			return nil, err
		}
	} else {
		// Создаем нового пользователя
		userID := uuid.New()

		user = &model.UserCredential{
			UserID:     userID.String(),
			Email:      email,
			IsVerified: true,
		}

		err = uc.userRepo.Create(ctx, user)
		if err != nil {
			return nil, err
		}

		// Создаем новую запись в OAuthAccount
		oauthAccount = &model.OauthAccount{
			UserID:         userID.String(),
			Provider:       provider,
			ProviderUserID: providerUserID,
		}

		err = uc.oauthRepo.Create(ctx, oauthAccount)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}
