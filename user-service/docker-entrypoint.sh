#!/bin/bash
set -e

# Запускаем ScyllaDB в фоне
/docker-entrypoint.py --smp 1 --memory 750M &

SCYLLA_PID=$!

# Ждем, пока ScyllaDB запустится
until cqlsh -e "DESCRIBE KEYSPACES"; do
  echo "Waiting for ScyllaDB to start..."
  sleep 5
done

# Выполняем инициализационный скрипт
cqlsh -f /init.cql

# Ожидаем завершения процесса ScyllaDB
wait $SCYLLA_PID