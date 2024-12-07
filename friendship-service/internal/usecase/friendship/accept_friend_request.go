package friendship

import (
	"context"
	"errors"
	"time"

	"github.com/malytinKonstantin/go-messenger-mono/friendship-service/internal/repositories"
)

type AcceptFriendRequestUsecase interface {
	Execute(ctx context.Context, requestID string) error
}

type acceptFriendRequestUsecase struct {
	repo repositories.FriendRequestRepository
}

func NewAcceptFriendRequestUsecase(repo repositories.FriendRequestRepository) AcceptFriendRequestUsecase {
	return &acceptFriendRequestUsecase{repo: repo}
}

func (uc *acceptFriendRequestUsecase) Execute(ctx context.Context, requestID string) error {
	if requestID == "" {
		return errors.New("request ID cannot be empty")
	}
	updatedAt := time.Now().Unix()
	err := uc.repo.UpdateFriendRequestStatus(ctx, requestID, "accepted", updatedAt)
	if err != nil {
		return err
	}
	return nil
}
