# Общие сущности и связи

## Сущность User (пользователь)

### Атрибуты:
- `user_id` (UUID): Уникальный идентификатор.
- `email` (string): Электронная почта.
- `nickname` (string): Никнейм.
- `password_hash` (string): Хэш пароля. (только в Auth Service)
- `bio` (string): Информация о себе.
- `avatar_url` (string): Ссылка на аватарку.

### Связи:
- User ↔ Friendship ↔ User:
  - Отношения дружбы между пользователями.
- User ↔ Message ↔ User:
  - Отправка и получение сообщений.
- User ↔ Notification:
  - Получение уведомлений о событиях.