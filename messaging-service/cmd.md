docker build -t messaging-cassandra-image -f Dockerfile.cassandra .

docker run --name messaging-cassandra \
    -p 9042:9042 \
    -d messaging-cassandra-image