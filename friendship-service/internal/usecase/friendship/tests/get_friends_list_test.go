package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/malytinKonstantin/go-messenger-mono/friendship-service/internal/models"
	"github.com/malytinKonstantin/go-messenger-mono/friendship-service/internal/usecase/friendship"
	"github.com/malytinKonstantin/go-messenger-mono/friendship-service/internal/usecase/friendship/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetFriendsListUsecaseExecute(t *testing.T) {
	ctx := context.Background()
	userID := "test-user-id"

	friends := []*models.User{
		{
			UserID:   "friend-1",
			Nickname: "Friend One",
		},
		{
			UserID:   "friend-2",
			Nickname: "Friend Two",
		},
	}

	mockRepo := new(mocks.UserRepository)
	mockRepo.On("GetFriends", mock.Anything, userID).Return(friends, nil)

	usecase := friendship.NewGetFriendsListUsecase(mockRepo)
	result, err := usecase.Execute(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, friends, result)
	mockRepo.AssertExpectations(t)
}

func TestGetFriendsListUsecaseExecuteEmptyList(t *testing.T) {
	ctx := context.Background()
	userID := "test-user-id"

	mockRepo := new(mocks.UserRepository)
	mockRepo.On("GetFriends", mock.Anything, userID).Return([]*models.User{}, nil)

	usecase := friendship.NewGetFriendsListUsecase(mockRepo)
	result, err := usecase.Execute(ctx, userID)

	assert.NoError(t, err)
	assert.Empty(t, result)
	mockRepo.AssertExpectations(t)
}

func TestGetFriendsListUsecaseExecuteError(t *testing.T) {
	ctx := context.Background()
	userID := "test-user-id"

	mockRepo := new(mocks.UserRepository)
	mockRepo.On("GetFriends", mock.Anything, userID).Return(nil, errors.New("database error"))

	usecase := friendship.NewGetFriendsListUsecase(mockRepo)
	result, err := usecase.Execute(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetFriendsListUsecaseExecuteWithEmptyUserID(t *testing.T) {
	ctx := context.Background()
	userID := ""

	mockRepo := new(mocks.UserRepository)
	usecase := friendship.NewGetFriendsListUsecase(mockRepo)
	friends, err := usecase.Execute(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, friends)
	assert.Equal(t, "user ID cannot be empty", err.Error())
	mockRepo.AssertNotCalled(t, "GetFriends")
}

func TestGetFriendsListUsecaseExecuteWithCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	userID := "test-user-id"

	mockRepo := new(mocks.UserRepository)
	mockRepo.On("GetFriends", ctx, userID).Return(nil, context.Canceled)

	usecase := friendship.NewGetFriendsListUsecase(mockRepo)
	result, err := usecase.Execute(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, context.Canceled, err)
	mockRepo.AssertExpectations(t)
}
