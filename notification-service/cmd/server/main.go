package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/notification_service/v1"
	"github.com/spf13/viper"
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
	env := os.Getenv("ENV")
	if env != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Файл .env не найден, продолжаем без него")
		}
	}
	viper.AutomaticEnv()

	// Чтение настроек Cassandra из переменных окружения
	host := viper.GetString("CASSANDRA_HOST")
	portStr := viper.GetString("CASSANDRA_PORT")
	keyspace := viper.GetString("CASSANDRA_KEYSPACE")
	username := viper.GetString("CASSANDRA_USERNAME")
	password := viper.GetString("CASSANDRA_PASSWORD")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Некорректный номер порта Cassandra: %v", err)
	}

	// Настройка кластера Cassandra
	cluster := gocql.NewCluster(host)
	cluster.Port = port
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum

	// Аутентификация (если требуется)
	if username != "" && password != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: username,
			Password: password,
		}
	}

	// Повторяем попытки подключения к Cassandra
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		session, err := cluster.CreateSession()
		if err == nil {
			defer session.Close()
			fmt.Println("Успешное подключение к Cassandra")
			break
		}
		log.Printf("Проблемы с подключением к Cassandra, попытка %d из %d: %v", i+1, maxRetries, err)
		time.Sleep(5 * time.Second)
		if i == maxRetries-1 {
			log.Fatalf("Не удалось подключиться к Cassandra после %d попыток", maxRetries)
		}
	}

	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Messaging Service is healthy")
	})

	// Запуск GRPC сервера
	grpcPort := viper.GetString("GRPC_PORT")
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Не удалось запустить слушатель: %v", err)
	}
	s := grpc.NewServer()
	// pb.RegisterMessagingServiceServer(s, &server{producer: p})

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Не удалось запустить GRPC сервер: %v", err)
		}
	}()

	// Запуск HTTP сервера
	httpPort := viper.GetString("HTTP_PORT")
	log.Fatal(app.Listen(":" + httpPort))
}
