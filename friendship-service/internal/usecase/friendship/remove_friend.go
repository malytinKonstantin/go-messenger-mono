package friendship

import (
	"context"
	"errors"

	"github.com/malytinKonstantin/go-messenger-mono/friendship-service/internal/repositories"
)

type RemoveFriendUsecase interface {
	Execute(ctx context.Context, userID, friendID string) error
}

type removeFriendUsecase struct {
	repo repositories.FriendRequestRepository
}

func NewRemoveFriendUsecase(repo repositories.FriendRequestRepository) RemoveFriendUsecase {
	return &removeFriendUsecase{repo: repo}
}

func (uc *removeFriendUsecase) Execute(ctx context.Context, userID, friendID string) error {
	if userID == "" {
		return errors.New("user ID cannot be empty")
	}
	if friendID == "" {
		return errors.New("friend ID cannot be empty")
	}

	err := uc.repo.DeleteFriendship(ctx, userID, friendID)
	if err != nil {
		return err
	}
	return nil
}
