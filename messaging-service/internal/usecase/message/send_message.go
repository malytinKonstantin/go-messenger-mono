package message

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/malytinKonstantin/go-messenger-mono/messaging-service/internal/models"
	"github.com/malytinKonstantin/go-messenger-mono/messaging-service/internal/repositories"
)

type SendMessageUsecase interface {
	Execute(ctx context.Context, message *models.Message) error
}

type sendMessageUsecase struct {
	messageRepo repositories.MessageRepository
	producer    *kafka.Producer
}

func NewSendMessageUsecase(messageRepo repositories.MessageRepository, producer *kafka.Producer) SendMessageUsecase {
	return &sendMessageUsecase{
		messageRepo: messageRepo,
		producer:    producer,
	}
}

func (uc *sendMessageUsecase) Execute(ctx context.Context, message *models.Message) error {
	return uc.messageRepo.SaveMessage(ctx, message)
}
