package user

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/malytinKonstantin/go-messenger-mono/user-service/internal/models"
	"github.com/malytinKonstantin/go-messenger-mono/user-service/internal/repositories"
)

type GetUserUsecase interface {
	Execute(ctx context.Context, userID gocql.UUID) (*models.UserProfile, error)
}

type getUserUsecase struct {
	userRepo repositories.UserRepository
}

func NewGetUserUsecase(userRepo repositories.UserRepository) GetUserUsecase {
	return &getUserUsecase{
		userRepo: userRepo,
	}
}

func (uc *getUserUsecase) Execute(ctx context.Context, userID gocql.UUID) (*models.UserProfile, error) {
	return uc.userRepo.GetUser(ctx, userID)
}
