package friendship

import (
	"context"

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

func (uc *getFriendsListUsecase) Execute(ctx context.Context, userID string) ([]*models.User, error) {
	friends, err := uc.repo.GetFriends(ctx, userID)
	if err != nil {
		return nil, err
	}
	return friends, nil
}
