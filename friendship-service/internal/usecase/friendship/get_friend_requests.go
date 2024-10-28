package friendship

import (
	"context"

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

func (uc *getFriendRequestsUsecase) Execute(ctx context.Context, userID string) ([]*models.FriendRequest, error) {
	requests, err := uc.repo.GetIncomingAndOutgoingRequests(ctx, userID)
	if err != nil {
		return nil, err
	}
	return requests, nil
}
