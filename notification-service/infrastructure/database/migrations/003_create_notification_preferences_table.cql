CREATE TABLE IF NOT EXISTS notification_service.notification_preferences (
    user_id uuid PRIMARY KEY,
    new_message boolean,
    friend_request boolean,
    system boolean
) WITH compaction = {'class': 'LeveledCompactionStrategy'}
  AND caching = {'keys': 'ALL', 'rows_per_partition': 'NONE'};