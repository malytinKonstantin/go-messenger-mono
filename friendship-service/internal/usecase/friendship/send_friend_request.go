package friendship

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/malytinKonstantin/go-messenger-mono/friendship-service/internal/models"
	"github.com/malytinKonstantin/go-messenger-mono/friendship-service/internal/repositories"
)

type SendFriendRequestUsecase interface {
	Execute(ctx context.Context, senderID, receiverID string) (*models.FriendRequest, error)
}

type sendFriendRequestUsecase struct {
	repo repositories.FriendRequestRepository
}

func NewSendFriendRequestUsecase(repo repositories.FriendRequestRepository) SendFriendRequestUsecase {
	return &sendFriendRequestUsecase{repo: repo}
}

func (uc *sendFriendRequestUsecase) Execute(ctx context.Context, senderID, receiverID string) (*models.FriendRequest, error) {
	if senderID == "" {
		return nil, errors.New("sender ID cannot be empty")
	}
	if receiverID == "" {
		return nil, errors.New("receiver ID cannot be empty")
	}
	if senderID == receiverID {
		return nil, errors.New("cannot send friend request to yourself")
	}

	request := &models.FriendRequest{
		RequestID:  uuid.New().String(),
		SenderID:   senderID,
		ReceiverID: receiverID,
		Status:     "pending",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := uc.repo.CreateFriendRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	return request, nil
}
