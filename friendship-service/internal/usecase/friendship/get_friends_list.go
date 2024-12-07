package friendship

import (
	"context"
	"errors"

	"github.com/malytinKonstantin/go-messenger-mono/friendship-service/internal/models"
	"github.com/malytinKonstantin/go-messenger-mono/friendship-service/internal/repositories"
)

type GetFriendsListUsecase interface {
	Execute(ctx context.Context, userID string) ([]*models.User, error)
}

type getFriendsListUsecase struct {
	repo repositories.UserRepository
}

func NewGetFriendsListUsecase(repo repositories.UserRepository) GetFriendsListUsecase {
	return &getFriendsListUsecase{repo: repo}
}

func (u *getFriendsListUsecase) Execute(ctx context.Context, userID string) ([]*models.User, error) {
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}

	return u.repo.GetFriends(ctx, userID)
}
