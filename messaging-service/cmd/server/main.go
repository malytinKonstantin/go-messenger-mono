package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

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

	// Прочая логика приложения
	// ...

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
