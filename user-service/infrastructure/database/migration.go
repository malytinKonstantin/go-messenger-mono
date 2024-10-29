package database

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

func RunMigrations(migrationsDir string) error {
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("cannot read migrations directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.HasSuffix(file.Name(), ".cql") {
			migrationPath := filepath.Join(migrationsDir, file.Name())
			content, err := ioutil.ReadFile(migrationPath)
			if err != nil {
				return fmt.Errorf("cannot read migration file %s: %v", file.Name(), err)
			}

			queries := strings.Split(string(content), ";")
			for _, query := range queries {
				query = strings.TrimSpace(query)
				if query == "" {
					continue
				}

				if err := ScyllaSession.Query(query).Exec(); err != nil {
					return fmt.Errorf("failed to execute migration %s: %v", file.Name(), err)
				}
				log.Printf("Executed migration: %s", file.Name())
			}
		}
	}
	return nil
}
