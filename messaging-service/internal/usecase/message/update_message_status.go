package message

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/malytinKonstantin/go-messenger-mono/messaging-service/internal/models"
	"github.com/malytinKonstantin/go-messenger-mono/messaging-service/internal/repositories"
)

type UpdateMessageStatusUsecase interface {
	Execute(ctx context.Context, messageID gocql.UUID, status models.MessageStatus) error
}

type updateMessageStatusUsecase struct {
	messageRepo repositories.MessageRepository
}

func NewUpdateMessageStatusUsecase(messageRepo repositories.MessageRepository) UpdateMessageStatusUsecase {
	return &updateMessageStatusUsecase{
		messageRepo: messageRepo,
	}
}

func (uc *updateMessageStatusUsecase) Execute(ctx context.Context, messageID gocql.UUID, status models.MessageStatus) error {
	return uc.messageRepo.UpdateMessageStatus(ctx, messageID, status)
}
