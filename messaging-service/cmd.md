docker build -t my-cassandra-image -f Dockerfile.cassandra .

docker run --name my-cassandra \
    -p 9042:9042 \
    -d my-cassandra-image