CREATE KEYSPACE IF NOT EXISTS mykeyspace WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

USE mykeyspace;

CREATE TABLE IF NOT EXISTS notifications (
    user_id UUID,
    notification_id TIMEUUID,
    type TEXT,
    content TEXT,
    is_read BOOLEAN,
    created_at TIMESTAMP,
    PRIMARY KEY (user_id, created_at, notification_id)
) WITH CLUSTERING ORDER BY (created_at DESC);

CREATE TABLE IF NOT EXISTS notification_preferences (
    user_id UUID PRIMARY KEY,
    preferences MAP<TEXT, BOOLEAN>
);