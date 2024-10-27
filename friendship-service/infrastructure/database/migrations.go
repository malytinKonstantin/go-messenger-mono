package database

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"

	"github.com/mindstand/gogm/v2"
)

type Migration struct {
	Version string
	Script  string
}

func RunMigrations() error {
	// Получаем сессию GOGM
	sess, err := Gogm.NewSessionV2(gogm.SessionConfig{AccessMode: gogm.Write})
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer sess.Close()

	// Убедимся, что существует узел для отслеживания миграций
	if err := createMigrationNode(sess); err != nil {
		return err
	}

	// Получаем список выполненных миграций
	appliedMigrations, err := getAppliedMigrations(sess)
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
			if err := applyMigration(sess, migration); err != nil {
				return err
			}
		}
	}

	return nil
}

func createMigrationNode(sess gogm.SessionV2) error {
	ctx := context.Background()
	query := `MERGE (m:MigrationVersion {id: 1})`
	_, err := sess.QueryRaw(ctx, query, nil)
	return err
}

func getAppliedMigrations(sess gogm.SessionV2) ([]string, error) {
	ctx := context.Background()
	var result struct {
		Versions []string `json:"versions"`
	}
	query := `MATCH (m:MigrationVersion) RETURN m.versions AS versions`
	records, err := sess.QueryRaw(ctx, query, nil)
	if err != nil {
		return nil, err
	}
	if len(records) > 0 {
		data, ok := records[0].Props["versions"].([]interface{})
		if ok {
			var versions []string
			for _, v := range data {
				versions = append(versions, v.(string))
			}
			return versions, nil
		}
	}
	return []string{}, nil
}

func applyMigration(sess gogm.SessionV2, migration Migration) error {
	ctx := context.Background()
	// Выполняем скрипт миграции
	_, err := sess.QueryRaw(ctx, migration.Script, nil)
	if err != nil {
		return fmt.Errorf("failed to apply migration %s: %v", migration.Version, err)
	}
	// Обновляем список выполненных миграций
	query := `
    MATCH (m:MigrationVersion {id: 1})
    SET m.versions = coalesce(m.versions, []) + $version
    `
	params := map[string]interface{}{
		"version": migration.Version,
	}
	_, err = sess.QueryRaw(ctx, query, params)
	return err
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
