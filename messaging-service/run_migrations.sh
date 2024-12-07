#!/bin/bash

CQLSH_HOST="${CASSANDRA_HOST:-localhost}"
CQLSH_PORT="${CASSANDRA_PORT:-9042}"

for file in $(ls migrations/*.cql | sort); do
  echo "Выполнение миграции: $file"
  cqlsh "$CQLSH_HOST" "$CQLSH_PORT" -f "$file"
done