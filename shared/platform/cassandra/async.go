package cassandra

import "fmt"

func (q *queryImpl) ExecAsync() (<-chan error, error) {
	errCh := make(chan error, 1)
	if err := q.query.Exec(); err != nil {
		errCh <- fmt.Errorf("async query execution failed: %w", err)
	}
	close(errCh)
	return errCh, nil
}
