package cassandra

import (
	"context"
	"fmt"
	"strings"
)

func (s *sessionImpl) CreateKeyspace(ctx context.Context, keyspace string, replication map[string]interface{}) error {
	replicationStr, err := formatReplication(replication)
	if err != nil {
		return err
	}
	query := fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s WITH replication = %s`, keyspace, replicationStr)
	return s.Query(query).WithContext(ctx).Exec()
}

func (s *sessionImpl) DropKeyspace(ctx context.Context, keyspace string) error {
	query := fmt.Sprintf(`DROP KEYSPACE IF EXISTS %s`, keyspace)
	return s.Query(query).WithContext(ctx).Exec()
}

func formatReplication(replication map[string]interface{}) (string, error) {
	if len(replication) == 0 {
		return "", fmt.Errorf("replication settings cannot be empty")
	}
	var entries []string
	for k, v := range replication {
		switch val := v.(type) {
		case string:
			entries = append(entries, fmt.Sprintf(`'%s': '%s'`, k, val))
		case int:
			entries = append(entries, fmt.Sprintf(`'%s': %d`, k, val))
		default:
			return "", fmt.Errorf("unsupported replication value type for key '%s'", k)
		}
	}
	return fmt.Sprintf("{%s}", strings.Join(entries, ", ")), nil
}
