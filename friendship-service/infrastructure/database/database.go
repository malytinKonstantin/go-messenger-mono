package database

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/spf13/viper"
)

var Driver neo4j.DriverWithContext

func InitNeo4jDriver() error {
	uri := fmt.Sprintf("bolt://%s:%d", viper.GetString("DATABASE_HOST"), viper.GetInt("DATABASE_PORT"))
	auth := neo4j.BasicAuth(viper.GetString("NEO4J_USER"), viper.GetString("NEO4J_PASSWORD"), "")
	var err error
	Driver, err = neo4j.NewDriverWithContext(uri, auth)
	if err != nil {
		return err
	}

	// Проверяем соединение
	ctx := context.Background()
	err = Driver.VerifyConnectivity(ctx)
	if err != nil {
		return err
	}

	return nil
}

func CloseNeo4jDriver() error {
	if Driver != nil {
		ctx := context.Background()
		return Driver.Close(ctx)
	}
	return nil
}
