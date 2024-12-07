CREATE KEYSPACE IF NOT EXISTS notification_service WITH replication = {
    'class': 'NetworkTopologyStrategy',
    'datacenter1': 1
};