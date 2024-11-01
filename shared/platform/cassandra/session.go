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

	if cfg.Username != "" && cfg.Password != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: cfg.Username,
			Password: cfg.Password,
		}
	}

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

func (s *sessionImpl) Prepare(stmt string) (PreparedStatement, error) {
	prepared, err := s.session.Prepare(stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	return &preparedStatementImpl{prepared: prepared}, nil
}

func (s *sessionImpl) NewBatch(batchType BatchType) Batch {
	return &batchImpl{
		batch:   gocql.NewBatch(gocql.BatchType(batchType)),
		session: s.session,
	}
}

func (s *sessionImpl) SetLogger(logger Logger) {
	s.logger = logger
	s.session.SetLogger(logger)
}

func (s *sessionImpl) Close() {
	s.session.Close()
}
