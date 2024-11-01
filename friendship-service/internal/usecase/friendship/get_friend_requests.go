package friendship

import (
	"context"
	"errors"

	"github.com/malytinKonstantin/go-messenger-mono/friendship-service/internal/models"
	"github.com/malytinKonstantin/go-messenger-mono/friendship-service/internal/repositories"
)

type GetFriendRequestsUsecase interface {
	Execute(ctx context.Context, userID string) ([]*models.FriendRequest, error)
}

type getFriendRequestsUsecase struct {
	repo repositories.FriendRequestRepository
}

func NewGetFriendRequestsUsecase(repo repositories.FriendRequestRepository) GetFriendRequestsUsecase {
	return &getFriendRequestsUsecase{repo: repo}
}

func (u *getFriendRequestsUsecase) Execute(ctx context.Context, userID string) ([]*models.FriendRequest, error) {
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}

	return u.repo.GetIncomingAndOutgoingRequests(ctx, userID)
}
