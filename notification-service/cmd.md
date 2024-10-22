docker build -t notification-cassandra-image -f Dockerfile.cassandra .

docker run --name notification-cassandra \
    -p 9242:9042 \
    -d notification-cassandra-image