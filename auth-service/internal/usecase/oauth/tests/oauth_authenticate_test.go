package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/infrastructure/database/model"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/repository"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/usecase/oauth"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/usecase/oauth/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOAuthAuthenticate_ExistingAccount(t *testing.T) {
	ctx := context.Background()
	userRepo := new(mocks.UserCredentialsRepository)
	oauthRepo := new(mocks.OauthAccountRepository)

	provider := "google"
	providerUserID := "provider-user-id"
	email := "user@example.com"
	userID := uuid.New().String()

	oauthAccount := &model.OauthAccount{
		UserID:         userID,
		Provider:       provider,
		ProviderUserID: providerUserID,
	}

	user := &model.UserCredential{
		UserID: userID,
		Email:  email,
	}

	oauthRepo.On("GetByProvider", mock.Anything, provider, providerUserID).Return(oauthAccount, nil)
	userRepo.On("GetByID", mock.Anything, mock.Anything).Return(user, nil)

	usecase := oauth.NewOAuthAuthenticateUsecase(userRepo, oauthRepo)

	result, err := usecase.Execute(ctx, provider, providerUserID, email)

	assert.NoError(t, err)
	assert.Equal(t, user, result)

	oauthRepo.AssertCalled(t, "GetByProvider", mock.Anything, provider, providerUserID)
	userRepo.AssertCalled(t, "GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID"))
}

func TestOAuthAuthenticate_NewAccount(t *testing.T) {
	ctx := context.Background()
	userRepo := new(mocks.UserCredentialsRepository)
	oauthRepo := new(mocks.OauthAccountRepository)

	provider := "google"
	providerUserID := "new-provider-user-id"
	email := "newuser@example.com"
	userID := uuid.New().String()

	oauthRepo.On("GetByProvider", mock.Anything, provider, providerUserID).Return(nil, repository.ErrNotFound)
	userRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.UserCredential")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*model.UserCredential)
		arg.UserID = userID
	})
	oauthRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.OauthAccount")).Return(nil)

	usecase := oauth.NewOAuthAuthenticateUsecase(userRepo, oauthRepo)

	result, err := usecase.Execute(ctx, provider, providerUserID, email)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, email, result.Email)
	assert.True(t, result.IsVerified)

	oauthRepo.AssertCalled(t, "GetByProvider", mock.Anything, provider, providerUserID)
	userRepo.AssertCalled(t, "Create", mock.Anything, mock.AnythingOfType("*model.UserCredential"))
	oauthRepo.AssertCalled(t, "Create", mock.Anything, mock.AnythingOfType("*model.OauthAccount"))
}

func TestOAuthAuthenticate_ErrorGettingOAuthAccount(t *testing.T) {
	ctx := context.Background()
	userRepo := new(mocks.UserCredentialsRepository)
	oauthRepo := new(mocks.OauthAccountRepository)

	provider := "google"
	providerUserID := "error-provider-user-id"
	email := "user@example.com"

	expectedError := errors.New("database error")

	oauthRepo.On("GetByProvider", mock.Anything, provider, providerUserID).Return(nil, expectedError)

	usecase := oauth.NewOAuthAuthenticateUsecase(userRepo, oauthRepo)

	result, err := usecase.Execute(ctx, provider, providerUserID, email)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)

	oauthRepo.AssertCalled(t, "GetByProvider", mock.Anything, provider, providerUserID)
}
