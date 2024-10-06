package main

import (
	"context"
	"log"
	"net"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gofiber/fiber/v2"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/notification-service"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedNotificationServiceServer
	producer *kafka.Producer
}

func (s *server) SendNotification(ctx context.Context, in *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	// Здесь должна быть реализация отправки уведомления
	// Для примера, просто логируем уведомление
	log.Printf("Отправка уведомления: UserID: %s, Message: %s, Type: %v", in.UserId, in.Message, in.Type)

	// Здесь можно добавить логику отправки уведомления через Kafka

	return &pb.SendNotificationResponse{Success: true}, nil
}

func (s *server) GetNotifications(ctx context.Context, in *pb.GetNotificationsRequest) (*pb.GetNotificationsResponse, error) {
	// Здесь должна быть реализация получения уведомлений
	// Для примера, возвращаем фиктивные данные
	notifications := []*pb.Notification{
		{
			Id:        "1",
			UserId:    in.UserId,
			Message:   "Тестовое уведомление",
			Type:      pb.NotificationType_SYSTEM,
			Timestamp: 1234567890,
		},
	}

	return &pb.GetNotificationsResponse{Notifications: notifications}, nil
}

func main() {
	// Инициализация Kafka producer
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Fatalf("Не удалось создать producer: %s", err)
	}
	defer p.Close()

	// Настройка gRPC сервера
	lis, err := net.Listen("tcp", ":50055")
	if err != nil {
		log.Fatalf("Не удалось прослушать порт: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterNotificationServiceServer(s, &server{producer: p})

	// Настройка Fiber для HTTP эндпоинтов
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Сервис уведомлений работает")
	})

	// Запуск HTTP сервера в отдельной горутине
	go func() {
		log.Fatal(app.Listen(":3005"))
	}()

	// Запуск gRPC сервера
	log.Printf("gRPC сервер слушает на %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
