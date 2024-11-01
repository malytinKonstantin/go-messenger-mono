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

func TestAcceptFriendRequestUsecaseExecute(t *testing.T) {
	ctx := context.Background()
	requestID := "test-request-id"

	mockRepo := new(mocks.FriendRequestRepository)
	mockRepo.On("UpdateFriendRequestStatus", mock.Anything, requestID, "accepted", mock.AnythingOfType("int64")).Return(nil)

	usecase := friendship.NewAcceptFriendRequestUsecase(mockRepo)
	err := usecase.Execute(ctx, requestID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAcceptFriendRequestUsecaseExecuteError(t *testing.T) {
	ctx := context.Background()
	requestID := "test-request-id"

	mockRepo := new(mocks.FriendRequestRepository)
	mockRepo.On("UpdateFriendRequestStatus", mock.Anything, requestID, "accepted", mock.AnythingOfType("int64")).
		Return(errors.New("database error"))

	usecase := friendship.NewAcceptFriendRequestUsecase(mockRepo)
	err := usecase.Execute(ctx, requestID)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestAcceptFriendRequestUsecaseExecuteWithEmptyID(t *testing.T) {
	ctx := context.Background()
	requestID := ""

	// Для пустого ID мы не должны вызывать репозиторий
	mockRepo := new(mocks.FriendRequestRepository)
	// Не настраиваем ожидания, так как метод не должен вызываться

	usecase := friendship.NewAcceptFriendRequestUsecase(mockRepo)
	err := usecase.Execute(ctx, requestID)

	assert.Error(t, err)
	assert.Equal(t, "request ID cannot be empty", err.Error())
	mockRepo.AssertNotCalled(t, "UpdateFriendRequestStatus")
}

func TestAcceptFriendRequestUsecaseExecuteWithCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	requestID := "test-request-id"

	mockRepo := new(mocks.FriendRequestRepository)
	mockRepo.On("UpdateFriendRequestStatus", ctx, requestID, "accepted", mock.AnythingOfType("int64")).
		Return(context.Canceled)

	usecase := friendship.NewAcceptFriendRequestUsecase(mockRepo)
	err := usecase.Execute(ctx, requestID)

	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
	mockRepo.AssertExpectations(t)
}
