# User Service (Сервис пользователей)

## Описание доменной области

User Service управляет профилями пользователей и связанной с ними информацией. Он предоставляет следующие возможности:

- Редактирование профиля:
  - Изменение никнейма (должен быть уникальным)
  - Обновление информации о себе
  - Загрузка аватарки
- Поиск пользователей по никнейму
- Хранение персональных данных

## Обоснование выбора базы данных: MongoDB

### Преимущества использования MongoDB:

1. **Гибкость схемы:** 
   - MongoDB - документоориентированная база данных
   - Позволяет хранить данные пользователей в гибком формате JSON
   - Легко изменять структуру без сложных миграций

2. **Высокая производительность:** 
   - Оптимизирована для быстрого чтения и записи документов
   - Улучшает отклик приложения при работе с профилями пользователей

3. **Масштабирование по горизонтали:** 
   - Поддержка шардинга
   - Эффективное распределение данных и нагрузки между серверами
   - Обеспечивает масштабируемость под высокие нагрузки

4. **Удобство работы с неструктурированными данными:** 
   - Легко хранить дополнительные поля в профиле пользователя
   - Не требует изменения схемы при добавлении новых атрибутов