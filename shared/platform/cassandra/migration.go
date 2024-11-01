package cassandra

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

func (s *sessionImpl) RunMigrations(ctx context.Context, migrationsDir string) error {
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrationFiles []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".cql" {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}

	sort.Strings(migrationFiles)

	for _, fileName := range migrationFiles {
		content, err := ioutil.ReadFile(filepath.Join(migrationsDir, fileName))
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", fileName, err)
		}

		queries := splitCQLStatements(string(content))

		for _, query := range queries {
			if err := s.Query(query).WithContext(ctx).Exec(); err != nil {
				return fmt.Errorf("failed to execute migration %s: %w", fileName, err)
			}
		}
	}

	return nil
}

func splitCQLStatements(content string) []string {
	return strings.Split(content, ";")
}
