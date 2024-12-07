package cassandra

import (
	"context"

	"github.com/gocql/gocql"
)

// batchImpl реализует интерфейс Batch.
type batchImpl struct {
	batch   *gocql.Batch
	session *gocql.Session
}

func (b *batchImpl) Query(stmt string, values ...interface{}) Batch {
	b.batch.Query(stmt, values...)
	return b
}

func (b *batchImpl) Exec() error {
	return b.session.ExecuteBatch(b.batch)
}

func (b *batchImpl) WithContext(ctx context.Context) Batch {
	b.batch = b.batch.WithContext(ctx)
	return b
}
