// Создание уникального ограничения для User.user_id
CREATE CONSTRAINT IF NOT EXISTS FOR (u:User)
REQUIRE u.user_id IS UNIQUE;

// Создание уникального ограничения для FriendRequest.request_id
CREATE CONSTRAINT IF NOT EXISTS FOR (fr:FriendRequest)
REQUIRE fr.request_id IS UNIQUE;

// Создание индекса для User.nickname
CREATE INDEX IF NOT EXISTS FOR (u:User)
ON (u.nickname);

// Создание индекса для User.avatar_url
CREATE INDEX IF NOT EXISTS FOR (u:User)
ON (u.avatar_url);

// Создание индексов для ускорения поиска запросов дружбы
CREATE INDEX IF NOT EXISTS FOR ()-[r:FRIEND_REQUEST]-()
ON (r.status);

CREATE INDEX IF NOT EXISTS FOR ()-[r:FRIEND_REQUEST]-()
ON (r.created_at);

// Создание ограничения для допустимых значений статуса запроса дружбы
CREATE CONSTRAINT IF NOT EXISTS ON ()-[r:FRIEND_REQUEST]-()
ASSERT r.status IN ['pending', 'accepted', 'rejected'];
