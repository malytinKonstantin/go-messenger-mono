#!/bin/sh

# Ждем, пока база данных будет готова
until PGPASSWORD=$DATABASE_PASSWORD psql -h $DATABASE_HOST -U $DATABASE_USER -d $DATABASE_NAME -c '\q'; do
  echo "Postgres is unavailable - sleeping"
  sleep 1
done

echo "Postgres is up - executing schema"

# Выполняем SQL-скрипты для создания схемы
PGPASSWORD=$DATABASE_PASSWORD psql -h $DATABASE_HOST -U $DATABASE_USER -d $DATABASE_NAME -f /app/db/schema/auth.sql

# Запускаем основное приложение
exec "$@"