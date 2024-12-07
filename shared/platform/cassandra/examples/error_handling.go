//go:build ignore

package examples

func error_handling() {
	// Пример обработки ошибок
	err := query.Exec()
	if err != nil {
		if errors.Is(err, cassandra.ErrNotFound) {
			log.Println("Record not found")
		} else {
			log.Printf("Query execution error: %v", err)
		}
	}
}
