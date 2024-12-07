CREATE KEYSPACE IF NOT EXISTS messaging_service WITH replication = {
    'class': 'NetworkTopologyStrategy',
    'datacenter1': 1
};