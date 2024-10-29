package database

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/spf13/viper"
)

var ScyllaSession *gocql.Session

func ConnectToScylla() error {
	host := viper.GetString("SCYLLA_HOST")
	port := viper.GetInt("SCYLLA_PORT")

	cluster := gocql.NewCluster(host)
	cluster.Port = port
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 10 * time.Second
	cluster.ConnectTimeout = 20 * time.Second
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: 5}
	cluster.DisableInitialHostLookup = true
	cluster.IgnorePeerAddr = true
	// Уберите аутентификацию, если она не нужна
	// cluster.Authenticator = gocql.PasswordAuthenticator{
	// 	Username: viper.GetString("SCYLLA_USERNAME"),
	// 	Password: viper.GetString("SCYLLA_PASSWORD"),
	// }

	session, err := cluster.CreateSession()
	if err != nil {
		return fmt.Errorf("error connecting to ScyllaDB: %w", err)
	}
	defer session.Close()

	// Создаем keyspace
	keyspace := viper.GetString("SCYLLA_KEYSPACE")
	if err := createKeyspace(session, keyspace); err != nil {
		return fmt.Errorf("error creating keyspace: %w", err)
	}

	// Переподключаемся с указанным keyspace
	cluster.Keyspace = keyspace
	ScyllaSession, err = cluster.CreateSession()
	if err != nil {
		return fmt.Errorf("error reconnecting to ScyllaDB with keyspace: %w", err)
	}

	return nil
}

func createKeyspace(session *gocql.Session, keyspace string) error {
	query := fmt.Sprintf(`
        CREATE KEYSPACE IF NOT EXISTS %s 
        WITH replication = {
            'class': 'SimpleStrategy',
            'replication_factor': '1'
        }`, keyspace)

	return session.Query(query).Exec()
}
