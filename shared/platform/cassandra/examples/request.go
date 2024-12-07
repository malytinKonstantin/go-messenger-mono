//go:build ignore

package examples

// Выполнение простого запроса на вставку данных

func test_request() {
	err := session.Query("INSERT INTO messages (id, content) VALUES (?, ?)", id, content).Exec()
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
	}
}
