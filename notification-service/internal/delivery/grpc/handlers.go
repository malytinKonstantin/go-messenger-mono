package grpc

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/malytinKonstantin/go-messenger-mono/notification-service/internal/models"
	"github.com/malytinKonstantin/go-messenger-mono/notification-service/internal/usecase/notification"
	"github.com/malytinKonstantin/go-messenger-mono/notification-service/internal/usecase/preferences"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/notification_service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NotificationServiceServer struct {
	pb.UnimplementedNotificationServiceServer
	sendNotificationUsecase       notification.SendNotificationUsecase
	getNotificationsUsecase       notification.GetNotificationsUsecase
	markNotificationAsReadUsecase notification.MarkNotificationAsReadUsecase
	updatePreferencesUsecase      preferences.UpdatePreferencesUsecase
	getPreferencesUsecase         preferences.GetPreferencesUsecase
}

func NewNotificationServiceServer(
	sendNotificationUC notification.SendNotificationUsecase,
	getNotificationsUC notification.GetNotificationsUsecase,
	markAsReadUC notification.MarkNotificationAsReadUsecase,
	updatePreferencesUC preferences.UpdatePreferencesUsecase,
	getPreferencesUC preferences.GetPreferencesUsecase,
) *NotificationServiceServer {
	return &NotificationServiceServer{
		sendNotificationUsecase:       sendNotificationUC,
		getNotificationsUsecase:       getNotificationsUC,
		markNotificationAsReadUsecase: markAsReadUC,
		updatePreferencesUsecase:      updatePreferencesUC,
		getPreferencesUsecase:         getPreferencesUC,
	}
}

// Отправка уведомления пользователю
func (s *NotificationServiceServer) SendNotification(ctx context.Context, req *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	userID, err := gocql.ParseUUID(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid user_id: %v", err)
	}

	notifType := models.NotificationType(req.Type)

	err = s.sendNotificationUsecase.Execute(ctx, userID, req.Message, notifType)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to send notification: %v", err)
	}

	return &pb.SendNotificationResponse{
		Success: true,
	}, nil
}

// Получение списка уведомлений пользователя
func (s *NotificationServiceServer) GetNotifications(ctx context.Context, req *pb.GetNotificationsRequest) (*pb.GetNotificationsResponse, error) {
	userID, err := gocql.ParseUUID(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid user_id: %v", err)
	}

	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10 // Значение по умолчанию
	}

	offset := int(req.Offset)
	if offset < 0 {
		offset = 0
	}

	// Получаем уведомления и общее количество
	notifications, totalCount, err := s.getNotificationsUsecase.Execute(ctx, userID, limit, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get notifications: %v", err)
	}

	// Преобразуем уведомления в формат protobuf
	pbNotifications := make([]*pb.Notification, len(notifications))
	for i, n := range notifications {
		pbNotifications[i] = &pb.Notification{
			Id:        n.NotificationID.String(),
			UserId:    n.UserID.String(),
			Message:   n.Message,
			Type:      pb.NotificationType(n.Type),
			CreatedAt: n.CreatedAt.Unix(),
			IsRead:    n.IsRead,
		}
	}

	// Формируем ответ
	resp := &pb.GetNotificationsResponse{
		Notifications: pbNotifications,
		TotalCount:    int32(totalCount),
	}

	return resp, nil
}

// Пометка уведомления как прочитанного
func (s *NotificationServiceServer) MarkNotificationAsRead(ctx context.Context, req *pb.MarkNotificationAsReadRequest) (*pb.MarkNotificationAsReadResponse, error) {
	notificationID, err := gocql.ParseUUID(req.NotificationId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid notification_id: %v", err)
	}

	userID, err := gocql.ParseUUID(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid user_id: %v", err)
	}

	err = s.markNotificationAsReadUsecase.Execute(ctx, notificationID, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to mark notification as read: %v", err)
	}

	return &pb.MarkNotificationAsReadResponse{
		Success: true,
	}, nil
}

// Обновление предпочтений уведомлений пользователя
func (s *NotificationServiceServer) UpdateNotificationPreferences(ctx context.Context, req *pb.UpdateNotificationPreferencesRequest) (*pb.UpdateNotificationPreferencesResponse, error) {
	userID, err := gocql.ParseUUID(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid user_id: %v", err)
	}

	preferences := &models.NotificationPreferences{
		UserID:        userID,
		NewMessage:    req.Preferences.NewMessage,
		FriendRequest: req.Preferences.FriendRequest,
		System:        req.Preferences.System,
	}

	err = s.updatePreferencesUsecase.Execute(ctx, preferences)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update notification preferences: %v", err)
	}

	return &pb.UpdateNotificationPreferencesResponse{
		Success: true,
	}, nil
}

// Получение предпочтений уведомлений пользователя
func (s *NotificationServiceServer) GetNotificationPreferences(ctx context.Context, req *pb.GetNotificationPreferencesRequest) (*pb.GetNotificationPreferencesResponse, error) {
	userID, err := gocql.ParseUUID(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid user_id: %v", err)
	}

	preferences, err := s.getPreferencesUsecase.Execute(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get notification preferences: %v", err)
	}

	pbPreferences := &pb.NotificationPreferences{
		NewMessage:    preferences.NewMessage,
		FriendRequest: preferences.FriendRequest,
		System:        preferences.System,
	}

	return &pb.GetNotificationPreferencesResponse{
		Preferences: pbPreferences,
	}, nil
}
