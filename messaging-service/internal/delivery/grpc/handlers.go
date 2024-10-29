package handlers

import (
	"context"
	"crypto/sha1"
	"time"

	"github.com/gocql/gocql"
	"github.com/malytinKonstantin/go-messenger-mono/messaging-service/internal/models"
	"github.com/malytinKonstantin/go-messenger-mono/messaging-service/internal/usecase/message"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/messaging_service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MessagingHandler struct {
	pb.UnimplementedMessagingServiceServer
	sendMessageUsecase         message.SendMessageUsecase
	getMessagesUsecase         message.GetMessagesUsecase
	updateMessageStatusUsecase message.UpdateMessageStatusUsecase
}

func NewMessagingHandler(
	sendMsgUc message.SendMessageUsecase,
	getMsgUc message.GetMessagesUsecase,
	updStatusUc message.UpdateMessageStatusUsecase,
) *MessagingHandler {
	return &MessagingHandler{
		sendMessageUsecase:         sendMsgUc,
		getMessagesUsecase:         getMsgUc,
		updateMessageStatusUsecase: updStatusUc,
	}
}

// Отправка сообщения
func (h *MessagingHandler) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	// Валидация запроса
	if req.SenderId == "" || req.RecipientId == "" || req.Content == "" {
		return nil, status.Errorf(codes.InvalidArgument, "sender_id, recipient_id and content must not be empty")
	}

	// Преобразование ID в UUID
	senderID, err := gocql.ParseUUID(req.SenderId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid sender_id: %v", err)
	}

	recipientID, err := gocql.ParseUUID(req.RecipientId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid recipient_id: %v", err)
	}

	// Определение ID беседы
	conversationID, err := generateConversationID(senderID, recipientID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error generating conversationID: %v", err)
	}

	// Создание сущности сообщения
	msg := &models.Message{
		MessageID:      gocql.TimeUUID(),
		SenderID:       senderID,
		RecipientID:    recipientID,
		ConversationID: conversationID,
		Content:        req.Content,
		Timestamp:      time.Now(),
		Status:         models.StatusSent,
	}

	// Запуск usecase отправки сообщения
	err = h.sendMessageUsecase.Execute(ctx, msg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error sending message: %v", err)
	}

	// Формирование ответа
	return &pb.SendMessageResponse{
		MessageId: msg.MessageID.String(),
	}, nil
}

// Получение сообщений
func (h *MessagingHandler) GetMessages(ctx context.Context, req *pb.GetMessagesRequest) (*pb.GetMessagesResponse, error) {
	// Валидация запроса
	if req.UserId == "" || req.ConversationUserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id and conversation_user_id must not be empty")
	}

	// Преобразование ID в UUID
	userID, err := gocql.ParseUUID(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	conversationUserID, err := gocql.ParseUUID(req.ConversationUserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid conversation_user_id: %v", err)
	}

	// Определение ID беседы
	conversationID, err := generateConversationID(userID, conversationUserID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error generating conversationID: %v", err)
	}

	// Запуск usecase получения сообщений
	messages, err := h.getMessagesUsecase.Execute(ctx, conversationID, int(req.Limit))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error getting messages: %v", err)
	}

	// Преобразование сообщений в формат proto
	pbMessages := make([]*pb.Message, len(messages))
	for i, msg := range messages {
		pbMessages[i] = &pb.Message{
			MessageId:   msg.MessageID.String(),
			SenderId:    msg.SenderID.String(),
			RecipientId: msg.RecipientID.String(),
			Content:     msg.Content,
			Timestamp:   msg.Timestamp.Unix(),
			Status:      pb.MessageStatus(msg.Status),
		}
	}

	// Формирование ответа
	return &pb.GetMessagesResponse{
		Messages: pbMessages,
	}, nil
}

// Обновление статуса сообщения
func (h *MessagingHandler) UpdateMessageStatus(ctx context.Context, req *pb.UpdateMessageStatusRequest) (*pb.UpdateMessageStatusResponse, error) {
	// Валидация запроса
	if req.MessageId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "message_id must not be empty")
	}

	// Преобразование ID в UUID
	messageID, err := gocql.ParseUUID(req.MessageId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid message_id: %v", err)
	}

	// Преобразование статуса
	statusMsg := models.MessageStatus(req.Status)

	// Обновление статуса сообщения
	err = h.updateMessageStatusUsecase.Execute(ctx, messageID, statusMsg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error updating message status: %v", err)
	}

	// Формирование ответа
	return &pb.UpdateMessageStatusResponse{
		Success: true,
	}, nil
}

// Генерация ID беседы на основе ID пользователей
func generateConversationID(user1ID, user2ID gocql.UUID) (gocql.UUID, error) {
	var (
		minID, maxID string
	)
	if user1ID.String() < user2ID.String() {
		minID = user1ID.String()
		maxID = user2ID.String()
	} else {
		minID = user2ID.String()
		maxID = user1ID.String()
	}

	h := sha1.New()
	h.Write([]byte(minID + maxID))
	hashed := h.Sum(nil)
	conversationID, err := gocql.UUIDFromBytes(hashed[:16])
	if err != nil {
		return gocql.UUID{}, err
	}
	return conversationID, nil
}
