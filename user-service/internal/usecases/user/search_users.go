package user

import (
	"context"

	"github.com/malytinKonstantin/go-messenger-mono/user-service/internal/models"
	"github.com/malytinKonstantin/go-messenger-mono/user-service/internal/repositories"
)

type SearchUsersUsecase interface {
	Execute(ctx context.Context, query string, limit, offset int) ([]*models.UserProfile, error)
}

type searchUsersUsecase struct {
	userRepo repositories.UserRepository
}

func NewSearchUsersUsecase(userRepo repositories.UserRepository) SearchUsersUsecase {
	return &searchUsersUsecase{
		userRepo: userRepo,
	}
}

func (uc *searchUsersUsecase) Execute(ctx context.Context, query string, limit, offset int) ([]*models.UserProfile, error) {
	return uc.userRepo.SearchUsers(ctx, query, limit, offset)
}
