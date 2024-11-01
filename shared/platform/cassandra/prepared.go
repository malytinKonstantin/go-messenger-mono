package cassandra

import (
	"github.com/gocql/gocql"
)

// preparedStatementImpl реализует интерфейс PreparedStatement.
type preparedStatementImpl struct {
	prepared *gocql.Prepared
}

func (p *preparedStatementImpl) Bind(values ...interface{}) Query {
	return &queryImpl{
		query: p.prepared.Bind(values...),
	}
}
