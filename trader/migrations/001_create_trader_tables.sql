-- migrations/001_create_trader_tables.sql
-- Run against the shared orderbook PostgreSQL database.
-- Tables are prefixed with "trader_" to avoid collision with orderbook tables.

CREATE TABLE IF NOT EXISTS trader_users (
    id            UUID        PRIMARY KEY,
    full_name     VARCHAR(150) NOT NULL,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT         NOT NULL,
    is_active     BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_trader_users_email ON trader_users (email);

-- ─────────────────────────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS trader_orders (
    id            UUID         PRIMARY KEY,
    user_id       UUID         NOT NULL REFERENCES trader_users (id) ON DELETE CASCADE,
    symbol        VARCHAR(20)  NOT NULL,
    side          VARCHAR(4)   NOT NULL CHECK (side IN ('bid', 'ask')),
    type          VARCHAR(10)  NOT NULL CHECK (type IN ('limit', 'market')),
    price         NUMERIC(30, 10),
    quantity      NUMERIC(30, 10) NOT NULL,
    status        VARCHAR(20)  NOT NULL DEFAULT 'submitted',
    raw_response  JSONB,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
ALTER TABLE trader_orders ADD COLUMN IF NOT EXISTS order_code VARCHAR(100);
CREATE INDEX IF NOT EXISTS idx_trader_orders_order_code ON trader_orders (order_code);
CREATE INDEX IF NOT EXISTS idx_trader_orders_user_id ON trader_orders (user_id);
CREATE INDEX IF NOT EXISTS idx_trader_orders_symbol  ON trader_orders (symbol);
CREATE INDEX IF NOT EXISTS idx_trader_orders_created ON trader_orders (created_at DESC);
