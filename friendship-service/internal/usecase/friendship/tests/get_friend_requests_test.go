package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/malytinKonstantin/go-messenger-mono/friendship-service/internal/models"
	"github.com/malytinKonstantin/go-messenger-mono/friendship-service/internal/usecase/friendship"
	"github.com/malytinKonstantin/go-messenger-mono/friendship-service/internal/usecase/friendship/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetFriendRequestsUsecaseExecute(t *testing.T) {
	ctx := context.Background()
	userID := "test-user-id"

	testRequests := []*models.FriendRequest{
		{
			RequestID:  "request-1",
			SenderID:   "user-1",
			ReceiverID: userID,
			Status:     "pending",
		},
		{
			RequestID:  "request-2",
			SenderID:   userID,
			ReceiverID: "user-2",
			Status:     "pending",
		},
	}

	mockRepo := new(mocks.FriendRequestRepository)
	mockRepo.On("GetIncomingAndOutgoingRequests", mock.Anything, userID).Return(testRequests, nil)

	usecase := friendship.NewGetFriendRequestsUsecase(mockRepo)
	requests, err := usecase.Execute(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, requests, 2)
	assert.Equal(t, testRequests, requests)
	mockRepo.AssertExpectations(t)
}

func TestGetFriendRequestsUsecaseExecuteEmptyList(t *testing.T) {
	ctx := context.Background()
	userID := "test-user-id"

	mockRepo := new(mocks.FriendRequestRepository)
	mockRepo.On("GetIncomingAndOutgoingRequests", mock.Anything, userID).Return([]*models.FriendRequest{}, nil)

	usecase := friendship.NewGetFriendRequestsUsecase(mockRepo)
	requests, err := usecase.Execute(ctx, userID)

	assert.NoError(t, err)
	assert.Empty(t, requests)
	mockRepo.AssertExpectations(t)
}

func TestGetFriendRequestsUsecaseExecuteError(t *testing.T) {
	ctx := context.Background()
	userID := "test-user-id"

	mockRepo := new(mocks.FriendRequestRepository)
	mockRepo.On("GetIncomingAndOutgoingRequests", mock.Anything, userID).
		Return(nil, errors.New("database error"))

	usecase := friendship.NewGetFriendRequestsUsecase(mockRepo)
	requests, err := usecase.Execute(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, requests)
	assert.Equal(t, "database error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetFriendRequestsUsecaseExecuteWithEmptyUserID(t *testing.T) {
	ctx := context.Background()
	userID := ""

	mockRepo := new(mocks.FriendRequestRepository)
	usecase := friendship.NewGetFriendRequestsUsecase(mockRepo)
	requests, err := usecase.Execute(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, requests)
	assert.Equal(t, "user ID cannot be empty", err.Error())
	mockRepo.AssertNotCalled(t, "GetIncomingAndOutgoingRequests")
}

func TestGetFriendRequestsUsecaseExecuteWithCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	userID := "test-user-id"

	mockRepo := new(mocks.FriendRequestRepository)
	mockRepo.On("GetIncomingAndOutgoingRequests", ctx, userID).Return(nil, context.Canceled)

	usecase := friendship.NewGetFriendRequestsUsecase(mockRepo)
	requests, err := usecase.Execute(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, requests)
	assert.Equal(t, context.Canceled, err)
	mockRepo.AssertExpectations(t)
}

func TestGetFriendRequestsUsecaseExecuteWithTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	time.Sleep(time.Millisecond * 2)
	userID := "test-user-id"

	mockRepo := new(mocks.FriendRequestRepository)
	mockRepo.On("GetIncomingAndOutgoingRequests", ctx, userID).Return(nil, context.DeadlineExceeded)

	usecase := friendship.NewGetFriendRequestsUsecase(mockRepo)
	requests, err := usecase.Execute(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, requests)
	assert.Equal(t, context.DeadlineExceeded, err)
	mockRepo.AssertExpectations(t)
}
