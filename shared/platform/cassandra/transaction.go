package cassandra

import (
	"context"
	"fmt"
)

func (q *queryImpl) WithSerialConsistency(consistency gocql.SerialConsistency) Query {
	q.query = q.query.SerialConsistency(consistency)
	return q
}

func (q *queryImpl) ExecCAS() (applied bool, err error) {
	applied, err = q.query.MapScanCAS(nil)
	if err != nil {
		return false, fmt.Errorf("CAS execution failed: %w", err)
	}
	return applied, nil
}
