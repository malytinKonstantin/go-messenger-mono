package database

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"

	"github.com/gocql/gocql"
)

func RunMigrations(session *gocql.Session, migrationsDir string) error {
	// Создаем таблицу для отслеживания миграций
	err := session.Query(`CREATE TABLE IF NOT EXISTS schema_migrations (
        version text PRIMARY KEY
    )`).Exec()
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %v", err)
	}

	// Получаем список уже выполненных миграций
	performedMigrations := map[string]bool{}
	iter := session.Query(`SELECT version FROM schema_migrations`).Iter()
	var version string
	for iter.Scan(&version) {
		performedMigrations[version] = true
	}
	if err := iter.Close(); err != nil {
		return fmt.Errorf("error reading performed migrations: %v", err)
	}

	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %v", err)
	}

	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && (file.Mode()&0400 != 0) {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}

	sort.Strings(migrationFiles)

	for _, fileName := range migrationFiles {
		if performedMigrations[fileName] {
			log.Printf("migration %s already performed, skipping...", fileName)
			continue
		}

		filePath := fmt.Sprintf("%s/%s", migrationsDir, fileName)
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %v", fileName, err)
		}

		// Разделяем содержимое файла на отдельные запросы
		queries := splitCQLStatements(string(content))
		log.Printf("executing migration: %s", fileName)

		for _, query := range queries {
			query = strings.TrimSpace(query)
			if query == "" {
				continue
			}
			log.Printf("executing query: %s", query)
			if err := session.Query(query).Exec(); err != nil {
				return fmt.Errorf("error executing migration %s: %v", fileName, err)
			}
		}

		// Сохранение информации о выполненной миграции
		if err := session.Query(`INSERT INTO schema_migrations (version) VALUES (?)`, fileName).Exec(); err != nil {
			return fmt.Errorf("failed to save migration version %s: %v", fileName, err)
		}
	}

	return nil
}

// Функция для разделения CQL-запросов по точке с запятой, учитывая комментарии и строки
func splitCQLStatements(content string) []string {
	var statements []string
	var sb strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "--") || strings.HasPrefix(trimmedLine, "//") || strings.HasPrefix(trimmedLine, "#") {
			// Пропускаем комментарии
			continue
		}
		sb.WriteString(line)
		if strings.HasSuffix(trimmedLine, ";") {
			statements = append(statements, sb.String())
			sb.Reset()
		} else {
			sb.WriteString(" ")
		}
	}
	// Добавляем последний запрос, если он не заканчивается на ';'
	if sb.Len() > 0 {
		statements = append(statements, sb.String())
	}
	return statements
}
