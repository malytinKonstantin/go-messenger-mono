package database

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/spf13/viper"
)

func ConnectToCassandra() (*gocql.Session, error) {
	cluster := gocql.NewCluster(viper.GetString("CASSANDRA_HOST"))
	cluster.Port = viper.GetInt("CASSANDRA_PORT")
	cluster.Keyspace = viper.GetString("CASSANDRA_KEYSPACE")
	cluster.Consistency = gocql.Quorum
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: viper.GetString("CASSANDRA_USERNAME"),
		Password: viper.GetString("CASSANDRA_PASSWORD"),
	}
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create Cassandra session: %v", err)
	}
	return session, nil
}
