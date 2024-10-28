package database

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/spf13/viper"
)

var Neo4jDriver neo4j.DriverWithContext

func InitNeo4jDriver() error {
	uri := viper.GetString("DATABASE_HOST")
	auth := neo4j.BasicAuth(viper.GetString("NEO4J_USER"), viper.GetString("NEO4J_PASSWORD"), "")
	var err error
	Neo4jDriver, err = neo4j.NewDriverWithContext(uri, auth)
	if err != nil {
		return err
	}

	// Проверяем соединение
	ctx := context.Background()
	err = Neo4jDriver.VerifyConnectivity(ctx)
	if err != nil {
		return err
	}

	return nil
}

func CloseNeo4jDriver() error {
	if Neo4jDriver != nil {
		ctx := context.Background()
		return Neo4jDriver.Close(ctx)
	}
	return nil
}
