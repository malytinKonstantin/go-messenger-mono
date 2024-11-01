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

func TestRejectFriendRequestUsecaseExecute(t *testing.T) {
	ctx := context.Background()
	requestID := "test-request-id"

	mockRepo := new(mocks.FriendRequestRepository)
	mockRepo.On("UpdateFriendRequestStatus", mock.Anything, requestID, "rejected", mock.AnythingOfType("int64")).Return(nil)

	usecase := friendship.NewRejectFriendRequestUsecase(mockRepo)
	err := usecase.Execute(ctx, requestID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRejectFriendRequestUsecaseExecuteError(t *testing.T) {
	ctx := context.Background()
	requestID := "test-request-id"

	mockRepo := new(mocks.FriendRequestRepository)
	mockRepo.On("UpdateFriendRequestStatus", mock.Anything, requestID, "rejected", mock.AnythingOfType("int64")).
		Return(errors.New("database error"))

	usecase := friendship.NewRejectFriendRequestUsecase(mockRepo)
	err := usecase.Execute(ctx, requestID)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestRejectFriendRequestUsecaseExecuteWithEmptyID(t *testing.T) {
	ctx := context.Background()
	requestID := ""

	mockRepo := new(mocks.FriendRequestRepository)
	usecase := friendship.NewRejectFriendRequestUsecase(mockRepo)
	err := usecase.Execute(ctx, requestID)

	assert.Error(t, err)
	assert.Equal(t, "request ID cannot be empty", err.Error())
	mockRepo.AssertNotCalled(t, "UpdateFriendRequestStatus")
}

func TestRejectFriendRequestUsecaseExecuteWithCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	requestID := "test-request-id"

	mockRepo := new(mocks.FriendRequestRepository)
	mockRepo.On("UpdateFriendRequestStatus", ctx, requestID, "rejected", mock.AnythingOfType("int64")).
		Return(context.Canceled)

	usecase := friendship.NewRejectFriendRequestUsecase(mockRepo)
	err := usecase.Execute(ctx, requestID)

	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
	mockRepo.AssertExpectations(t)
}
