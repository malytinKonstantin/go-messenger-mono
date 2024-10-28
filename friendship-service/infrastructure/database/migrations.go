package database

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Migration struct {
	Version string
	Script  string
}

func RunMigrations() error {
	ctx := context.Background()
	// Создаем сессию
	session := Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	// Убедимся, что существует узел для отслеживания миграций
	if err := createMigrationNode(ctx, session); err != nil {
		return err
	}

	// Получаем список выполненных миграций
	appliedMigrations, err := getAppliedMigrations(ctx, session)
	if err != nil {
		return err
	}

	// Читаем файлы миграций из директории
	migrationFiles, err := ioutil.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %v", err)
	}

	// Сортируем файлы миграций
	var migrations []Migration
	for _, file := range migrationFiles {
		if strings.HasSuffix(file.Name(), ".cql") {
			version := strings.Split(file.Name(), "_")[0]
			script, err := ioutil.ReadFile("migrations/" + file.Name())
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %v", file.Name(), err)
			}
			migrations = append(migrations, Migration{
				Version: version,
				Script:  string(script),
			})
		}
	}
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	// Применяем новые миграции
	for _, migration := range migrations {
		if !contains(appliedMigrations, migration.Version) {
			log.Printf("Applying migration %s", migration.Version)
			if err := applyMigration(ctx, session, migration); err != nil {
				return err
			}
		}
	}

	return nil
}

func createMigrationNode(ctx context.Context, session neo4j.SessionWithContext) error {
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `MERGE (m:MigrationVersion {id: 1})`
		_, err := tx.Run(ctx, query, nil)
		return nil, err
	})
	return err
}

func getAppliedMigrations(ctx context.Context, session neo4j.SessionWithContext) ([]string, error) {
	var appliedMigrations []string
	_, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `MATCH (m:MigrationVersion) RETURN m.versions AS versions`
		result, err := tx.Run(ctx, query, nil)
		if err != nil {
			return nil, err
		}
		if result.Next(ctx) {
			record := result.Record()
			if versions, ok := record.Get("versions"); ok && versions != nil {
				if vs, ok := versions.([]interface{}); ok {
					for _, v := range vs {
						appliedMigrations = append(appliedMigrations, v.(string))
					}
				}
			}
		}
		return nil, result.Err()
	})
	return appliedMigrations, err
}

func applyMigration(ctx context.Context, session neo4j.SessionWithContext, migration Migration) error {
	// Выполняем скрипт миграции
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(ctx, migration.Script, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to apply migration %s: %v", migration.Version, err)
		}

		// Обновляем список выполненных миграций
		query := `
			MATCH (m:MigrationVersion {id: 1})
			SET m.versions = coalesce(m.versions, []) + $version
		`
		params := map[string]interface{}{
			"version": migration.Version,
		}
		_, err = tx.Run(ctx, query, params)
		return nil, err
	})
	return err
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
