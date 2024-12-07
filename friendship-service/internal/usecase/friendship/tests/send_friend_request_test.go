package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/malytinKonstantin/go-messenger-mono/friendship-service/internal/usecase/friendship"
	"github.com/malytinKonstantin/go-messenger-mono/friendship-service/internal/usecase/friendship/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSendFriendRequestUsecaseExecute(t *testing.T) {
	ctx := context.Background()
	senderID := uuid.New().String()
	receiverID := uuid.New().String()

	mockRepo := new(mocks.FriendRequestRepository)
	mockRepo.On("CreateFriendRequest", mock.Anything, mock.AnythingOfType("*models.FriendRequest")).Return(nil)

	usecase := friendship.NewSendFriendRequestUsecase(mockRepo)
	request, err := usecase.Execute(ctx, senderID, receiverID)

	assert.NoError(t, err)
	assert.NotNil(t, request)
	assert.Equal(t, senderID, request.SenderID)
	assert.Equal(t, receiverID, request.ReceiverID)
	assert.Equal(t, "pending", request.Status)
	assert.WithinDuration(t, time.Now(), request.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), request.UpdatedAt, time.Second)
	mockRepo.AssertExpectations(t)
}

func TestSendFriendRequestUsecaseExecuteError(t *testing.T) {
	ctx := context.Background()
	senderID := uuid.New().String()
	receiverID := uuid.New().String()

	mockRepo := new(mocks.FriendRequestRepository)
	mockRepo.On("CreateFriendRequest", mock.Anything, mock.AnythingOfType("*models.FriendRequest")).
		Return(errors.New("database error"))

	usecase := friendship.NewSendFriendRequestUsecase(mockRepo)
	request, err := usecase.Execute(ctx, senderID, receiverID)

	assert.Error(t, err)
	assert.Nil(t, request)
	assert.Equal(t, "database error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestSendFriendRequestUsecaseExecuteWithEmptyIDs(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name       string
		senderID   string
		receiverID string
	}{
		{"empty sender ID", "", "receiver-id"},
		{"empty receiver ID", "sender-id", ""},
		{"both IDs empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.FriendRequestRepository)
			usecase := friendship.NewSendFriendRequestUsecase(mockRepo)
			request, err := usecase.Execute(ctx, tt.senderID, tt.receiverID)

			assert.Error(t, err)
			assert.Nil(t, request)
			assert.Contains(t, err.Error(), "cannot be empty")
		})
	}
}

func TestSendFriendRequestUsecaseExecuteWithCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	senderID := uuid.New().String()
	receiverID := uuid.New().String()

	mockRepo := new(mocks.FriendRequestRepository)
	mockRepo.On("CreateFriendRequest", ctx, mock.AnythingOfType("*models.FriendRequest")).
		Return(context.Canceled)

	usecase := friendship.NewSendFriendRequestUsecase(mockRepo)
	request, err := usecase.Execute(ctx, senderID, receiverID)

	assert.Error(t, err)
	assert.Nil(t, request)
	assert.Equal(t, context.Canceled, err)
	mockRepo.AssertExpectations(t)
}

func TestSendFriendRequestUsecaseSameUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New().String()

	mockRepo := new(mocks.FriendRequestRepository)
	usecase := friendship.NewSendFriendRequestUsecase(mockRepo)
	request, err := usecase.Execute(ctx, userID, userID)

	assert.Error(t, err)
	assert.Nil(t, request)
	assert.Equal(t, "cannot send friend request to yourself", err.Error())
}
