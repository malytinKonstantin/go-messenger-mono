CREATE TABLE IF NOT EXISTS notification_service.notifications (
    user_id uuid,
    created_at timestamp,
    notification_id timeuuid,
    message text,
    type int,
    is_read boolean,
    PRIMARY KEY ((user_id), created_at, notification_id)
) WITH CLUSTERING ORDER BY (created_at DESC, notification_id DESC)
  AND compaction = {'class': 'TimeWindowCompactionStrategy'}
  AND default_time_to_live = 0
  AND gc_grace_seconds = 864000;