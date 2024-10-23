package database

import (
	"fmt"
	"strconv"

	"github.com/gocql/gocql"
	"github.com/spf13/viper"
)

func ConnectToScylla() (*gocql.Session, error) {
	host := viper.GetString("SCYLLA_HOST")
	port := viper.GetString("SCYLLA_PORT")
	keyspace := viper.GetString("SCYLLA_KEYSPACE")
	consistency := viper.GetString("SCYLLA_CONSISTENCY")
	username := viper.GetString("SCYLLA_USERNAME")
	password := viper.GetString("SCYLLA_PASSWORD")

	portNum, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("invalid port number: %v", err)
	}

	cluster := gocql.NewCluster(host)
	cluster.Port = portNum
	cluster.Keyspace = keyspace

	switch consistency {
	case "ONE":
		cluster.Consistency = gocql.One
	case "QUORUM":
		cluster.Consistency = gocql.Quorum
	case "ALL":
		cluster.Consistency = gocql.All
	default:
		cluster.Consistency = gocql.Quorum
	}

	if username != "" && password != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: username,
			Password: password,
		}
	}

	cluster.IgnorePeerAddr = true
	cluster.DisableInitialHostLookup = true

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("error connecting to ScyllaDB: %v", err)
	}
	return session, nil
}
