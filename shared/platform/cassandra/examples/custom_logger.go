//go:build ignore

package examples

import "log"

type CustomLogger struct{}

func (l *CustomLogger) Info(msg string) {
	log.Printf("INFO: %s", msg)
}

func (l *CustomLogger) Error(msg string) {
	log.Printf("ERROR: %s", msg)
}

func custom_logger() {
	// ...
	session.SetLogger(&CustomLogger{})
	// Теперь сессия будет использовать ваш логгер для сообщений
}
