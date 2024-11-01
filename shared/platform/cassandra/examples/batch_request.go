//go:build ignore

package examples

func batch_request() {
	// Запрос с пагинацией
	query := session.Query("SELECT * FROM messages").PageSize(20)
	iter := query.Iter()
	var messageID, content string
	for iter.Scan(&messageID, &content) {
		// Обработка каждой записи
		log.Printf("Message ID: %s, Content: %s", messageID, content)
	}
	if err := iter.Close(); err != nil {
		log.Printf("Failed to close iterator: %v", err)
	}

	// Получение PageState для следующей страницы
	pageState := iter.PageState()

	// Следующий запрос с использованием PageState
	nextQuery := session.Query("SELECT * FROM messages").PageSize(20).PageState(pageState)
	nextIter := nextQuery.Iter()
	// ... повторите процесс

}
