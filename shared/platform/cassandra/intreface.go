package cassandra

import (
	"context"
	"errors"
	"time"

	"github.com/gocql/gocql"
)

// Config представляет настройки подключения к Cassandra.
type Config struct {
	Hosts          []string
	Port           int
	Keyspace       string
	Consistency    gocql.Consistency
	ConnectTimeout time.Duration
	Username       string
	Password       string
	MaxOpenConns   int // Максимальное количество открытых соединений
	// Добавьте другие необходимые параметры (например, SSL настройки)
}

// Session представляет сессию соединения с базой данных Cassandra.
type Session interface {
	// Query создает новый запрос.
	Query(stmt string, values ...interface{}) Query
	// Close закрывает сессию.
	Close()
	// Prepare подготавливает выражение для последующего выполнения.
	Prepare(stmt string) (PreparedStatement, error)
	// NewBatch создает новый батч-запрос.
	NewBatch(batchType BatchType) Batch
	// SetLogger устанавливает логгер для сессии.
	SetLogger(logger Logger)
}

// Query представляет запрос к базе данных Cassandra.
type Query interface {
	// Exec выполняет запрос без возвращения результатов.
	Exec() error
	// Scan выполняет запрос и сканирует единственную строку результата.
	Scan(dest ...interface{}) error
	// Iter выполняет запрос и возвращает итератор для обхода результатов.
	Iter() Iterator
	// WithContext связывает контекст с запросом.
	WithContext(ctx context.Context) Query
	// PageSize устанавливает размер страницы для результатов запроса.
	PageSize(n int) Query
	// PageState устанавливает состояние страницы для запросов пагинации.
	PageState(state []byte) Query
	// Consistency устанавливает уровень консистентности для запроса.
	Consistency(consistency gocql.Consistency) Query
	// Bind связывает параметры с подготовленным запросом.
	Bind(values ...interface{}) Query
}

// Iterator представляет итератор для обхода результатов запроса.
type Iterator interface {
	// Scan считывает следующую строку результатов.
	Scan(dest ...interface{}) bool
	// Close закрывает итератор.
	Close() error
	// PageState возвращает состояние страницы для продолжения пагинации.
	PageState() []byte
	// WillSwitchPage указывает, произойдет ли переключение страницы.
	WillSwitchPage() bool
	// Host возвращает информацию о хосте, с которым установлено соединение.
	Host() *gocql.HostInfo
}

// BatchType определяет тип батч-запроса.
type BatchType int

const (
	LoggedBatch BatchType = iota
	UnloggedBatch
	CounterBatch
)

// Batch представляет батч-запрос к базе данных Cassandra.
type Batch interface {
	// Query добавляет запрос в батч.
	Query(stmt string, values ...interface{}) Batch
	// Exec выполняет батч-запрос.
	Exec() error
	// WithContext связывает контекст с батчем.
	WithContext(ctx context.Context) Batch
}

// PreparedStatement представляет подготовленное выражение.
type PreparedStatement interface {
	// Bind связывает параметры с подготовленным запросом.
	Bind(values ...interface{}) Query
}

// Logger представляет интерфейс для логирования.
type Logger interface {
	// Info логирует информационное сообщение.
	Info(msg string)
	// Error логирует сообщение об ошибке.
	Error(msg string)
}

// Определение ошибок.
var (
	ErrNotFound        = gocql.ErrNotFound
	ErrConnection      = gocql.ErrNoConnections
	ErrQueryExecution  = errors.New("query execution error")
	ErrInvalidArgument = gocql.ErrInvalid
)
