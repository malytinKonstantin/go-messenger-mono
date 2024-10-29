package message

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/malytinKonstantin/go-messenger-mono/messaging-service/internal/models"
	"github.com/malytinKonstantin/go-messenger-mono/messaging-service/internal/repositories"
)

type GetMessagesUsecase interface {
	Execute(ctx context.Context, conversationID gocql.UUID, limit int) ([]*models.Message, error)
}

type getMessagesUsecase struct {
	messageRepo repositories.MessageRepository
}

func NewGetMessagesUsecase(messageRepo repositories.MessageRepository) GetMessagesUsecase {
	return &getMessagesUsecase{
		messageRepo: messageRepo,
	}
}

func (uc *getMessagesUsecase) Execute(ctx context.Context, conversationID gocql.UUID, limit int) ([]*models.Message, error) {
	return uc.messageRepo.GetMessages(ctx, conversationID, limit)
}
