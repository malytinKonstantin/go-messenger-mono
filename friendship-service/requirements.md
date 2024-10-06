# Сервис дружбы (Friendship Service)

## Основные функции

### Управление друзьями
- Отправка запроса на добавление в друзья
- Принятие или отклонение входящего запроса
- Удаление пользователя из друзей

### Просмотр списка друзей
- Получение списка друзей
- Получение списка входящих и исходящих запросов

## Функциональные требования

### Обработка отношений
- Поддержка состояний дружбы (ожидает, подтверждено, отклонено)

### Уведомления
- Интеграция с Notification Service для оповещений о новых запросах и действиях

## Сущности и атрибуты

### Friendship
- `user_id` (UUID): Идентификатор пользователя
- `friend_id` (UUID): Идентификатор друга
- `status` (enum): Статус отношения ("pending", "accepted", "rejected")
- `created_at` (datetime): Дата отправки запроса
- `updated_at` (datetime): Дата последнего обновления статуса