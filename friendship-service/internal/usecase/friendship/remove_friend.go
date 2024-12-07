package friendship

import (
	"context"

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
	err := uc.repo.DeleteFriendship(ctx, userID, friendID)
	if err != nil {
		return err
	}
	return nil
}
