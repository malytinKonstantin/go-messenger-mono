package database

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/malytinKonstantin/go-messenger-mono/shared/platform/cassandra"
)

func DbConfig() cassandra.Config {
	cassandraConfig := cassandra.Config{
		Hosts:          []string{viper.GetString("CASSANDRA_HOST")},
		Port:           viper.GetInt("CASSANDRA_PORT"),
		Keyspace:       viper.GetString("CASSANDRA_KEYSPACE"),
		Consistency:    cassandra.ParseConsistency(viper.GetString("CASSANDRA_CONSISTENCY")),
		ConnectTimeout: viper.GetDuration("CASSANDRA_CONNECT_TIMEOUT"),
		Username:       viper.GetString("CASSANDRA_USERNAME"),
		Password:       viper.GetString("CASSANDRA_PASSWORD"),
		MaxOpenConns:   viper.GetInt("CASSANDRA_MAX_OPEN_CONNS"),
	}
	return cassandraConfig
}

func ConnectToCassandra() (cassandra.Session, error) {
	cfg := DbConfig()
	session, err := cassandra.NewSession(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Cassandra: %w", err)
	}
	return session, nil
}

func CreateKeyspace(session cassandra.Session) error {
	keyspaceQuery := `CREATE KEYSPACE IF NOT EXISTS notification_service WITH replication = {
        'class': 'NetworkTopologyStrategy',
        'datacenter1': 1
    };`
	if err := session.Query(keyspaceQuery).Exec(); err != nil {
		return fmt.Errorf("failed to create keyspace: %v", err)
	}
	return nil
}

func ReconnectToCassandraWithKeyspace() (cassandra.Session, error) {
	cfg := DbConfig()
	cfg.Keyspace = "notification_service"
	session, err := cassandra.NewSession(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to reconnect to Cassandra with keyspace: %w", err)
	}
	return session, nil
}

func RunMigrations(session cassandra.Session, migrationsDir string) error {
	if err := cassandra.RunMigrations(session, migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}
