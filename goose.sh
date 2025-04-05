#!/bin/bash

DB_PATH="./payments.db"
MIGRATIONS_DIR="./storage/migrations"

case "$1" in
  up)
    goose -dir $MIGRATIONS_DIR sqlite3 $DB_PATH up
    ;;
  down)
    goose -dir $MIGRATIONS_DIR sqlite3 $DB_PATH down
    ;;
  create)
    if [ -z "$2" ]; then
      echo "Usage: ./goose.sh create migration_name"
      exit 1
    fi
    goose -dir $MIGRATIONS_DIR create $2 sql
    ;;
  *)
    echo "Usage: ./goose.sh {up|down|create migration_name}"
    exit 1
    ;;
esac