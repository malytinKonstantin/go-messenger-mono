package friendship

import (
	"context"
	"time"

	"github.com/malytinKonstantin/go-messenger-mono/friendship-service/internal/repositories"
)

type RejectFriendRequestUsecase interface {
	Execute(ctx context.Context, requestID string) error
}

type rejectFriendRequestUsecase struct {
	repo repositories.FriendRequestRepository
}

func NewRejectFriendRequestUsecase(repo repositories.FriendRequestRepository) RejectFriendRequestUsecase {
	return &rejectFriendRequestUsecase{repo: repo}
}

func (uc *rejectFriendRequestUsecase) Execute(ctx context.Context, requestID string) error {
	updatedAt := time.Now().Unix()
	err := uc.repo.UpdateFriendRequestStatus(ctx, requestID, "rejected", updatedAt)
	if err != nil {
		return err
	}
	return nil
}
