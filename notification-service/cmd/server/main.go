package main

import (
	"context"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/notification-service"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedNotificationServiceServer
	producer         *kafka.Producer
	cassandraSession *gocql.Session
}

func (s *server) SendNotification(ctx context.Context, in *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	// Генерация уникального идентификатора для уведомления
	notificationID := gocql.TimeUUID()
	createdAt := time.Now()

	// Сохранение уведомления в Cassandra
	err := s.cassandraSession.Query(`
        INSERT INTO notifications (user_id, notification_id, type, content, is_read, created_at)
        VALUES (?, ?, ?, ?, ?, ?)`,
		in.UserId, notificationID, in.Type.String(), in.Message, false, createdAt).Exec()
	if err != nil {
		log.Printf("Ошибка при сохранении уведомления в Cassandra: %v", err)
		return &pb.SendNotificationResponse{Success: false}, err
	}

	log.Printf("Уведомление сохранено: UserID: %s, Message: %s, Type: %v", in.UserId, in.Message, in.Type)

	// Дополнительная логика отправки уведомления

	return &pb.SendNotificationResponse{Success: true}, nil
}

func (s *server) GetNotifications(ctx context.Context, in *pb.GetNotificationsRequest) (*pb.GetNotificationsResponse, error) {
	// Получение уведомлений из Cassandra
	var notifications []*pb.Notification

	iter := s.cassandraSession.Query(`
        SELECT notification_id, type, content, is_read, created_at
        FROM notifications WHERE user_id = ?`, in.UserId).Iter()

	var notificationID gocql.UUID
	var typeStr string
	var content string
	var isRead bool
	var createdAt time.Time

	for iter.Scan(&notificationID, &typeStr, &content, &isRead, &createdAt) {
		notificationType := pb.NotificationType(pb.NotificationType_value[typeStr])

		notifications = append(notifications, &pb.Notification{
			Id:        notificationID.String(),
			UserId:    in.UserId,
			Message:   content,
			Type:      notificationType,
			CreatedAt: createdAt.Unix(),
			IsRead:    isRead,
		})
	}

	if err := iter.Close(); err != nil {
		log.Printf("Ошибка при получении уведомлений: %v", err)
		return nil, err
	}

	return &pb.GetNotificationsResponse{Notifications: notifications}, nil
}

func main() {
	// Инициализация Kafka producer
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka:9092"})
	if err != nil {
		log.Fatalf("Не удалось создать producer: %s", err)
	}
	defer p.Close()

	// Подключение к Cassandra
	cassandraHost := os.Getenv("CASSANDRA_HOST")
	cassandraPortStr := os.Getenv("CASSANDRA_PORT")
	cassandraKeyspace := os.Getenv("CASSANDRA_KEYSPACE")

	cassandraPort, err := strconv.Atoi(cassandraPortStr)
	if err != nil {
		log.Fatalf("Некорректный порт Cassandra: %v", err)
	}

	cluster := gocql.NewCluster(cassandraHost)
	cluster.Port = cassandraPort
	cluster.Keyspace = cassandraKeyspace
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Не удалось подключиться к Cassandra: %v", err)
	}
	defer session.Close()

	// Настройка gRPC сервера
	lis, err := net.Listen("tcp", ":50055")
	if err != nil {
		log.Fatalf("Не удалось прослушать порт: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterNotificationServiceServer(s, &server{
		producer:         p,
		cassandraSession: session,
	})

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
