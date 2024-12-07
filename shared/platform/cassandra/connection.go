package cassandra

import (
	"fmt"

	"github.com/gocql/gocql"
)

// sessionImpl реализует интерфейс Session.
type sessionImpl struct {
	session *gocql.Session
	logger  Logger
}

func NewSession(cfg Config) (Session, error) {
	cluster := gocql.NewCluster(cfg.Hosts...)
	cluster.Port = cfg.Port
	cluster.Keyspace = cfg.Keyspace
	cluster.Consistency = cfg.Consistency
	cluster.ConnectTimeout = cfg.ConnectTimeout

	// Аутентификация
	if cfg.Username != "" && cfg.Password != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: cfg.Username,
			Password: cfg.Password,
		}
	}

	// Настройка пулов соединений
	if cfg.MaxOpenConns > 0 {
		cluster.NumConns = cfg.MaxOpenConns
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create Cassandra session: %w", err)
	}

	return &sessionImpl{session: session}, nil
}

func (s *sessionImpl) Query(stmt string, values ...interface{}) Query {
	return &queryImpl{
		query: s.session.Query(stmt, values...),
	}
}

func (s *sessionImpl) Close() {
	s.session.Close()
}

func (s *sessionImpl) NewBatch(batchType gocql.BatchType) Batch {
	return &batchImpl{
		batch:   gocql.NewBatch(batchType),
		session: s.session,
	}
}

func (s *sessionImpl) Prepare(stmt string) (PreparedStatement, error) {
	prepared, err := s.session.Prepare(stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	return &preparedStatementImpl{prepared: prepared}, nil
}

func (s *sessionImpl) SetLogger(logger Logger) {
	s.logger = logger
}
