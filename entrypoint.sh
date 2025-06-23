#!/bin/sh
set -e

DB_FILE="/root/db/ontrek.db"
MIGRATIONS_DIR="/root/migrations"

sqlite3 "$DB_FILE" "CREATE TABLE IF NOT EXISTS migrations (filename TEXT PRIMARY KEY, applied_at DATETIME DEFAULT CURRENT_TIMESTAMP);"

for f in "$MIGRATIONS_DIR"/*.sql; do
  [ -e "$f" ] || continue

  FILENAME=$(basename "$f")
  if ! sqlite3 "$DB_FILE" "SELECT 1 FROM migrations WHERE filename = '$FILENAME'" | grep -q 1; then
    echo ">> Applying migration: $FILENAME"
    sqlite3 "$DB_FILE" < "$f"
    sqlite3 "$DB_FILE" "INSERT INTO migrations (filename) VALUES ('$FILENAME');"
  fi
done

echo ">> Starting server..."
exec ./server
