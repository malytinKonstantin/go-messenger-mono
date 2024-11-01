package cassandra

import "log"

// Пример простой реализации Logger.
type SimpleLogger struct{}

func (l *SimpleLogger) Info(msg string) {
	log.Printf("INFO: %s", msg)
}

func (l *SimpleLogger) Error(msg string) {
	log.Printf("ERROR: %s", msg)
}
