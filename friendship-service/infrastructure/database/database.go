package database

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/spf13/viper"
)

func ConnectToNeo4j() (neo4j.Driver, error) {
	driver, err := neo4j.NewDriver(
		fmt.Sprintf("bolt://%s:%s", viper.GetString("DATABASE_HOST"), viper.GetString("DATABASE_PORT")),
		neo4j.BasicAuth(
			viper.GetString("NEO4J_USER"),
			viper.GetString("NEO4J_PASSWORD"),
			"",
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Neo4j driver: %v", err)
	}
	return driver, nil
}
