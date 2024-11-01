package cassandra

import (
	"context"

	"github.com/gocql/gocql"
)

// queryImpl реализует интерфейс Query.
type queryImpl struct {
	query *gocql.Query
}

func (q *queryImpl) Exec() error {
	return q.query.Exec()
}

func (q *queryImpl) Scan(dest ...interface{}) error {
	return q.query.Scan(dest...)
}

func (q *queryImpl) Iter() Iterator {
	return &iteratorImpl{
		iter: q.query.Iter(),
	}
}

func (q *queryImpl) WithContext(ctx context.Context) Query {
	q.query.WithContext(ctx)
	return q
}

func (q *queryImpl) PageSize(n int) Query {
	q.query.PageSize(n)
	return q
}

func (q *queryImpl) PageState(state []byte) Query {
	q.query.PageState(state)
	return q
}

func (q *queryImpl) Consistency(consistency gocql.Consistency) Query {
	q.query.Consistency(consistency)
	return q
}

func (q *queryImpl) Bind(values ...interface{}) Query {
	q.query.Bind(values...)
	return q
}
