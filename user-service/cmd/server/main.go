package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/user_service/v1"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUserServiceServer
	producer *kafka.Producer
}

func (s *server) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// Здесь должна быть реализация получения пользователя из базы данных
	// Для примера возвращаем фиктивные данные
	profile := &pb.UserProfile{
		UserId:    in.UserId,
		Nickname:  "Иван Иванов",
		Bio:       "Пример биографии",
		AvatarUrl: "https://example.com/avatar.jpg",
	}

	// Отправляем событие в Kafka
	topic := "user_events"
	err := s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte("Получен пользователь: " + profile.UserId),
	}, nil)

	if err != nil {
		log.Printf("Не удалось отправить сообщение: %v", err)
	}

	return &pb.GetUserResponse{Profile: profile}, nil
}

func (s *server) CreateUserProfile(ctx context.Context, in *pb.CreateUserProfileRequest) (*pb.CreateUserProfileResponse, error) {
	// Здесь должна быть реализация создания профиля пользователя в базе данных
	// Для примера создаем фиктивный профиль пользователя
	profile := &pb.UserProfile{
		UserId:    in.UserId,
		Nickname:  in.Nickname,
		Bio:       in.Bio,
		AvatarUrl: in.AvatarUrl,
	}

	// Отправляем событие в Kafka
	topic := "user_events"
	err := s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte("Создан новый профиль пользователя: " + profile.UserId),
	}, nil)

	if err != nil {
		log.Printf("Не удалось отправить сообщение: %v", err)
	}

	return &pb.CreateUserProfileResponse{Profile: profile}, nil
}

func (s *server) UpdateUserProfile(ctx context.Context, in *pb.UpdateUserProfileRequest) (*pb.UpdateUserProfileResponse, error) {
	// Здесь должна быть реализация обновления профиля пользователя в базе данных
	// Для примера обновляем фиктивный профиль пользователя
	profile := &pb.UserProfile{
		UserId:    in.UserId,
		Nickname:  in.Nickname,
		Bio:       in.Bio,
		AvatarUrl: in.AvatarUrl,
	}

	// Отправляем событие в Kafka
	topic := "user_events"
	err := s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte("Обновлен профиль пользователя: " + profile.UserId),
	}, nil)

	if err != nil {
		log.Printf("Не удалось отправить сообщение: %v", err)
	}

	return &pb.UpdateUserProfileResponse{Profile: profile}, nil
}

func (s *server) SearchUsers(ctx context.Context, in *pb.SearchUsersRequest) (*pb.SearchUsersResponse, error) {
	// Здесь должна быть реализация поиска пользователей в базе данных
	// Для примера возвращаем список фиктивных профилей
	profiles := []*pb.UserProfile{
		{
			UserId:    "user1",
			Nickname:  "User One",
			Bio:       "Биография пользователя один",
			AvatarUrl: "https://example.com/avatar1.jpg",
		},
		{
			UserId:    "user2",
			Nickname:  "User Two",
			Bio:       "Биография пользователя два",
			AvatarUrl: "https://example.com/avatar2.jpg",
		},
	}

	// Отправляем событие в Kafka
	topic := "user_events"
	err := s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte("Произведен поиск пользователей по запросу: " + in.Query),
	}, nil)

	if err != nil {
		log.Printf("Не удалось отправить сообщение: %v", err)
	}

	return &pb.SearchUsersResponse{Profiles: profiles}, nil
}

func main() {
	env := os.Getenv("ENV")
	if env != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Файл .env не найден, продолжаем без него")
		}
	}

	// Автоматически считываем переменные окружения
	viper.AutomaticEnv()

	// Читаем параметры из переменных окружения
	grpcPort := viper.GetString("GRPC_PORT")
	httpPort := viper.GetString("HTTP_PORT")

	// Читаем параметры ScyllaDB из переменных окружения
	scyllaHost := viper.GetString("SCYLLA_HOST")
	scyllaPort := viper.GetString("SCYLLA_PORT")
	scyllaKeyspace := viper.GetString("SCYLLA_KEYSPACE")
	scyllaConsistency := viper.GetString("SCYLLA_CONSISTENCY")

	fmt.Println(scyllaHost, scyllaPort, scyllaKeyspace, scyllaConsistency)

	// Настраиваем подключение к ScyllaDB
	cluster := gocql.NewCluster(fmt.Sprintf("%s:%s", scyllaHost, scyllaPort))
	cluster.Keyspace = scyllaKeyspace

	// Устанавливаем уровень согласованности
	switch scyllaConsistency {
	case "ONE":
		cluster.Consistency = gocql.One
	case "QUORUM":
		cluster.Consistency = gocql.Quorum
	case "ALL":
		cluster.Consistency = gocql.All
	default:
		cluster.Consistency = gocql.Quorum
	}

	session, err := cluster.CreateSession()
	if err != nil {
		fmt.Printf("Ошибка подключения к ScyllaDB: %v\n", err)
		return
	}
	defer session.Close()

	// Настройка Kafka producer
	// p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	// if err != nil {
	// 	log.Fatalf("Не удалось создать producer: %s", err)
	// }
	// defer p.Close()

	// Запуск gRPC сервера
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("Не удалось прослушать порт: %v", err)
	}
	s := grpc.NewServer()
	// pb.RegisterUserServiceServer(s, &server{producer: p})

	// Запуск Fiber сервера для внутреннего API
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Сервис пользователей работает")
	})

	go func() {
		log.Fatal(app.Listen(fmt.Sprintf(":%s", httpPort)))
	}()

	log.Printf("gRPC сервер слушает на %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
