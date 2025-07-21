#!/bin/sh
set -e

BACKUP_DIR="/backups"
DB_FILE="/root/db/ontrek.db"
MIGRATIONS_DIR="/root/migrations"
mkdir -p "$BACKUP_DIR"

if [ -f "$DB_FILE" ]; then
  TIMESTAMP=$(date +"%Y%m%d-%H%M%S")
  BACKUP_FILE="$BACKUP_DIR/ontrek-$TIMESTAMP.db.bak"
  cp "$DB_FILE" "$BACKUP_FILE"
  echo ">> Backup created: $BACKUP_FILE"

# Mantieni solo gli ultimi 5 backup
BACKUP_PATTERN="$BACKUP_DIR/ontrek-*.db.bak"
BACKUPS_TO_DELETE=$(ls -1t $BACKUP_PATTERN | tail -n +6)
if [ -n "$BACKUPS_TO_DELETE" ]; then
  echo ">> Removing old backups..."
  echo "$BACKUPS_TO_DELETE" | xargs rm -f
fi
else
  echo ">> No one database file found, skipping backup."
fi

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
