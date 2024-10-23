package database

import (
	"fmt"
	"strconv"

	"github.com/gocql/gocql"
	"github.com/spf13/viper"
)

func ConnectToCassandra() (*gocql.Session, error) {
	host := viper.GetString("CASSANDRA_HOST")
	portStr := viper.GetString("CASSANDRA_PORT")
	keyspace := viper.GetString("CASSANDRA_KEYSPACE")
	username := viper.GetString("CASSANDRA_USERNAME")
	password := viper.GetString("CASSANDRA_PASSWORD")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid Cassandra port number: %v", err)
	}

	cluster := gocql.NewCluster(host)
	cluster.Port = port
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum

	if username != "" && password != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: username,
			Password: password,
		}
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Cassandra: %v", err)
	}
	return session, nil
}
