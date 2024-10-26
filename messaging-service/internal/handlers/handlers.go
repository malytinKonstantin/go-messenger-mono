package handlers

import (
	"context"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gocql/gocql"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/messaging_service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MessagingHandler struct {
	pb.UnimplementedMessagingServiceServer
	producer *kafka.Producer
	session  *gocql.Session
}

func NewMessagingHandler(producer *kafka.Producer, session *gocql.Session) *MessagingHandler {
	return &MessagingHandler{producer: producer, session: session}
}

func (h *MessagingHandler) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation error: %v", err)
	}

	// Мокаем ответ
	return &pb.SendMessageResponse{
		MessageId: "550e8400-e29b-41d4-a716-446655440000",
	}, nil
}

func (h *MessagingHandler) GetMessages(ctx context.Context, req *pb.GetMessagesRequest) (*pb.GetMessagesResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation error: %v", err)
	}

	// Мокаем список сообщений
	messages := []*pb.Message{
		{
			MessageId:   "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			SenderId:    "6ba7b811-9dad-11d1-80b4-00c04fd430c8",
			RecipientId: req.UserId,
			Content:     "Привет!",
			Timestamp:   time.Now().Unix(),
			Status:      pb.MessageStatus_MESSAGE_STATUS_SENT,
		},
		{
			MessageId:   "6ba7b812-9dad-11d1-80b4-00c04fd430c8",
			SenderId:    req.UserId,
			RecipientId: "6ba7b811-9dad-11d1-80b4-00c04fd430c8",
			Content:     "Как дела?",
			Timestamp:   time.Now().Unix(),
			Status:      pb.MessageStatus_MESSAGE_STATUS_READ,
		},
	}

	return &pb.GetMessagesResponse{
		Messages: messages,
	}, nil
}

func (h *MessagingHandler) UpdateMessageStatus(ctx context.Context, req *pb.UpdateMessageStatusRequest) (*pb.UpdateMessageStatusResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation error: %v", err)
	}

	// Мокаем успешный ответ
	return &pb.UpdateMessageStatusResponse{
		Success: true,
	}, nil
}
