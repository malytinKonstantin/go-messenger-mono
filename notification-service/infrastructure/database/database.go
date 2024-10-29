package database

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/spf13/viper"
)

func ConnectToCassandra() (*gocql.Session, error) {
	cluster := gocql.NewCluster(viper.GetString("CASSANDRA_HOST"))
	cluster.Port = viper.GetInt("CASSANDRA_PORT")
	cluster.Consistency = gocql.Quorum
	return cluster.CreateSession()
}

func ReconnectToCassandraWithKeyspace() (*gocql.Session, error) {
	cluster := gocql.NewCluster(viper.GetString("CASSANDRA_HOST"))
	cluster.Port = viper.GetInt("CASSANDRA_PORT")
	cluster.Keyspace = "notification_service"
	cluster.Consistency = gocql.Quorum
	return cluster.CreateSession()
}

func CreateKeyspace(session *gocql.Session) error {
	keyspaceQuery := `CREATE KEYSPACE IF NOT EXISTS notification_service WITH replication = {
        'class': 'NetworkTopologyStrategy',
        'datacenter1': 1
    };`
	if err := session.Query(keyspaceQuery).Exec(); err != nil {
		return fmt.Errorf("failed to create keyspace: %v", err)
	}
	return nil
}
