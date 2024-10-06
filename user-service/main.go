package main

import (
	"context"
	"log"
	"net"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gofiber/fiber/v2"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/user-service"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUserServiceServer
	producer *kafka.Producer
}

func (s *server) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.User, error) {
	// Здесь должна быть реализация получения пользователя из базы данных
	// Для примера возвращаем фиктивные данные
	user := &pb.User{Id: in.UserId, Name: "Иван Иванов", Email: "ivan@example.com"}

	// Отправляем событие в Kafka
	topic := "user_events"
	err := s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte("Получен пользователь: " + user.Id),
	}, nil)

	if err != nil {
		log.Printf("Не удалось отправить сообщение: %v", err)
	}

	return user, nil
}

func (s *server) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// Здесь должна быть реализация создания пользователя в базе данных
	// Для примера создаем фиктивного пользователя
	user := &pb.User{
		Id:    "new_user_id", // В реальном приложении это должен быть уникальный идентификатор
		Name:  in.Name,
		Email: in.Email,
	}

	// Отправляем событие в Kafka
	topic := "user_events"
	err := s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte("Создан новый пользователь: " + user.Id),
	}, nil)

	if err != nil {
		log.Printf("Не удалось отправить сообщение: %v", err)
	}

	return &pb.CreateUserResponse{User: user}, nil
}

func main() {
	// Настройка Kafka producer
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Fatalf("Не удалось создать producer: %s", err)
	}
	defer p.Close()

	// Запуск gRPC сервера
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Не удалось прослушать порт: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{producer: p})

	// Запуск Fiber сервера для внутреннего API
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Сервис пользователей работает")
	})

	go func() {
		log.Fatal(app.Listen(":3001"))
	}()

	log.Printf("gRPC сервер слушает на %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
