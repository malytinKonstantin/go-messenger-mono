//go:build ignore

package examples

func prepared_example() {
	// Подготовка выражения
	prepared, err := session.Prepare("INSERT INTO messages (id, content) VALUES (?, ?)")
	if err != nil {
		log.Printf("Failed to prepare statement: %v", err)
	}

	// Связывание параметров и выполнение
	query := prepared.Bind(id, content)
	if err := query.Exec(); err != nil {
		log.Printf("Failed to execute prepared statement: %v", err)
	}
}
