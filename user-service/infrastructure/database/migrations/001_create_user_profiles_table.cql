CREATE KEYSPACE IF NOT EXISTS user_service WITH replication = {
    'class': 'SimpleStrategy',
    'replication_factor': '1'
};

CREATE TABLE IF NOT EXISTS user_service.user_profiles (
    user_id uuid PRIMARY KEY,
    nickname text,
    bio text,
    avatar_url text,
    created_at timestamp,
    updated_at timestamp
);