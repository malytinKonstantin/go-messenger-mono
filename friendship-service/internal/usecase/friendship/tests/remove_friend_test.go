package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/malytinKonstantin/go-messenger-mono/friendship-service/internal/usecase/friendship"
	"github.com/malytinKonstantin/go-messenger-mono/friendship-service/internal/usecase/friendship/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRemoveFriendUsecaseExecute(t *testing.T) {
	ctx := context.Background()
	userID := "test-user-id"
	friendID := "friend-user-id"

	mockRepo := new(mocks.FriendRequestRepository)
	mockRepo.On("DeleteFriendship", mock.Anything, userID, friendID).Return(nil)

	usecase := friendship.NewRemoveFriendUsecase(mockRepo)
	err := usecase.Execute(ctx, userID, friendID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRemoveFriendUsecaseExecuteError(t *testing.T) {
	ctx := context.Background()
	userID := "test-user-id"
	friendID := "friend-user-id"

	mockRepo := new(mocks.FriendRequestRepository)
	mockRepo.On("DeleteFriendship", mock.Anything, userID, friendID).
		Return(errors.New("database error"))

	usecase := friendship.NewRemoveFriendUsecase(mockRepo)
	err := usecase.Execute(ctx, userID, friendID)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestRemoveFriendUsecaseExecuteWithEmptyIDs(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name     string
		userID   string
		friendID string
	}{
		{"empty user ID", "", "friend-id"},
		{"empty friend ID", "user-id", ""},
		{"both IDs empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.FriendRequestRepository)
			usecase := friendship.NewRemoveFriendUsecase(mockRepo)
			err := usecase.Execute(ctx, tt.userID, tt.friendID)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "cannot be empty")
		})
	}
}

func TestRemoveFriendUsecaseExecuteWithCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	userID := "test-user-id"
	friendID := "friend-user-id"

	mockRepo := new(mocks.FriendRequestRepository)
	mockRepo.On("DeleteFriendship", ctx, userID, friendID).Return(context.Canceled)

	usecase := friendship.NewRemoveFriendUsecase(mockRepo)
	err := usecase.Execute(ctx, userID, friendID)

	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
	mockRepo.AssertExpectations(t)
}
