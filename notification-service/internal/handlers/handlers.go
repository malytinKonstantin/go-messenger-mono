package handlers

import (
	"context"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/notification_service/v1"
)

type NotificationServiceServer struct {
	pb.UnimplementedNotificationServiceServer
	producer *kafka.Producer
	session  *gocql.Session
}

func NewNotificationServiceServer(producer *kafka.Producer, session *gocql.Session) *NotificationServiceServer {
	return &NotificationServiceServer{producer: producer, session: session}
}

func (s *NotificationServiceServer) SendNotification(ctx context.Context, req *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	// Моковая логика отправки уведомления
	return &pb.SendNotificationResponse{
		Success: true,
	}, nil
}

func (s *NotificationServiceServer) GetNotifications(ctx context.Context, req *pb.GetNotificationsRequest) (*pb.GetNotificationsResponse, error) {
	// Моковые данные уведомлений
	notifications := []*pb.Notification{
		{
			Id:        uuid.New().String(),
			UserId:    req.UserId,
			Message:   "У вас новое сообщение от пользователя Иван Иванов",
			Type:      pb.NotificationType_NOTIFICATION_TYPE_NEW_MESSAGE,
			CreatedAt: time.Now().Unix(),
			IsRead:    false,
		},
		{
			Id:        uuid.New().String(),
			UserId:    req.UserId,
			Message:   "Пользователь Мария Петрова отправила вам запрос в друзья",
			Type:      pb.NotificationType_NOTIFICATION_TYPE_FRIEND_REQUEST,
			CreatedAt: time.Now().Add(-2 * time.Hour).Unix(),
			IsRead:    false,
		},
		{
			Id:        uuid.New().String(),
			UserId:    req.UserId,
			Message:   "Системное уведомление: новое приложение доступно для загрузки",
			Type:      pb.NotificationType_NOTIFICATION_TYPE_SYSTEM,
			CreatedAt: time.Now().Add(-24 * time.Hour).Unix(),
			IsRead:    true,
		},
	}

	start := req.Offset
	end := start + req.Limit
	totalNotifications := int32(len(notifications))

	if start > totalNotifications {
		paginatedNotifications := []*pb.Notification{}

		return &pb.GetNotificationsResponse{
			Notifications: paginatedNotifications,
		}, nil
	}

	if end > totalNotifications || req.Limit == 0 {
		end = totalNotifications
	}

	paginatedNotifications := notifications[start:end]

	return &pb.GetNotificationsResponse{
		Notifications: paginatedNotifications,
	}, nil
}

func (s *NotificationServiceServer) MarkNotificationAsRead(ctx context.Context, req *pb.MarkNotificationAsReadRequest) (*pb.MarkNotificationAsReadResponse, error) {
	// Моковая логика пометки уведомления как прочитанного
	return &pb.MarkNotificationAsReadResponse{
		Success: true,
	}, nil
}

func (s *NotificationServiceServer) UpdateNotificationPreferences(ctx context.Context, req *pb.UpdateNotificationPreferencesRequest) (*pb.UpdateNotificationPreferencesResponse, error) {
	// Моковая логика обновления предпочтений

	return &pb.UpdateNotificationPreferencesResponse{
		Success: true,
	}, nil
}

func (s *NotificationServiceServer) GetNotificationPreferences(ctx context.Context, req *pb.GetNotificationPreferencesRequest) (*pb.GetNotificationPreferencesResponse, error) {
	// Моковые данные предпочтений уведомлений
	preferences := &pb.NotificationPreferences{
		NewMessage:    true,
		FriendRequest: true,
		System:        false,
	}

	return &pb.GetNotificationPreferencesResponse{
		Preferences: preferences,
	}, nil
}
