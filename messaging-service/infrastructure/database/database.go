package database

import (
	"github.com/gocql/gocql"
	"github.com/spf13/viper"
)

func ConnectToCassandra() (*gocql.Session, error) {
	cluster := gocql.NewCluster(viper.GetString("CASSANDRA_HOST"))
	cluster.Port = viper.GetInt("CASSANDRA_PORT")
	cluster.Keyspace = "messaging_service"
	cluster.Consistency = gocql.One

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}
