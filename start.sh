#!/bin/sh

set -e

echo "run db migration"
/app/migrate -path /app/migrations -database "$DB_SOURCE" -verbose up

echo "start the app"

#takes all parameters passed to the script and run it
exec "$@"