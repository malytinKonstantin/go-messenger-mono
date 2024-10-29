CREATE TABLE IF NOT EXISTS messages (
    conversation_id uuid,
    message_id timeuuid,
    sender_id uuid,
    recipient_id uuid,
    content text,
    status text,
    PRIMARY KEY (conversation_id, message_id)
) WITH CLUSTERING ORDER BY (message_id DESC)
  AND compaction = {'class': 'LeveledCompactionStrategy'}
  AND caching = {'keys': 'ALL', 'rows_per_partition': 'ALL'};