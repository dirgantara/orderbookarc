#!/bin/sh
set -e

echo "[trader] waiting for postgres..."
until pg_isready -h "${PGHOST:-postgres}" -U "${PGUSER:-orderbook}" -q; do
  sleep 1
done
echo "[trader] postgres is ready"

echo "[trader] running migrations..."
psql "${DATABASE_URL}" -f migrations/001_create_trader_tables.sql
echo "[trader] migrations done"

echo "[trader] starting server on :${SERVER_PORT:-9000}"

# 🔥 run langsung tanpa build
exec go run ./cmd/main.go