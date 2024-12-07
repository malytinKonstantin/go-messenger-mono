//go:build ignore

package examples

import (
	"log"
	"time"

	"github.com/gocql/gocql"
	"github.com/malytinKonstantin/go-messenger-mono/shared/platform/cassandra"
)

func main() {
	cfg := cassandra.Config{
		Hosts:          []string{"127.0.0.1"},
		Port:           9042,
		Keyspace:       "messaging_service",
		Consistency:    gocql.Quorum,
		ConnectTimeout: 10 * time.Second,
		Username:       "your_username",
		Password:       "your_password",
		MaxOpenConns:   10, // Устанавливаем максимальное количество открытых соединений
	}

	session, err := cassandra.NewSession(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}
	defer session.Close()

	// Использование сессии
	query := session.Query("SELECT * FROM messages WHERE user_id = ?", "some-user-id")
	iter := query.Iter()
	// Обработка результатов
	var messageID string
	var content string
	for iter.Scan(&messageID, &content) {
		log.Printf("Message ID: %s, Content: %s", messageID, content)
	}
	if err := iter.Close(); err != nil {
		log.Printf("Failed to close iterator: %v", err)
	}
}
