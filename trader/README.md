# Trader — Mini Client Service

Go service yang berjalan di port **9000**, terhubung ke postgres yang sama dengan orderbook, dan memanggil `POST /api/v1/orders/` milik orderbook menggunakan **HMAC authentication**.

## Struktur Project

```
trader/
├── cmd/main.go                      # Entrypoint
├── internal/
│   ├── config/config.go             # Env-based config
│   ├── handler/handler.go           # HTTP handlers (Gin)
│   ├── middleware/jwt.go            # JWT Bearer middleware
│   ├── model/model.go               # User & Order structs
│   ├── repository/repository.go     # Postgres queries (pgx)
│   └── service/
│       ├── auth.go                  # Register + Login + JWT sign
│       ├── order.go                 # PlaceOrder via orderbook API
│       └── hmac.go                  # HMAC signer (mirrors orderbook middleware)
├── migrations/
│   └── 001_create_trader_tables.sql # DDL — trader_users, trader_orders
├── Dockerfile
└── entrypoint.sh                    # Wait for PG → migrate → start
```

## API Endpoints

| Method | Path | Auth | Keterangan |
|--------|------|------|-----------|
| `POST` | `/auth/register` | — | Daftar user baru |
| `POST` | `/auth/login` | — | Login, dapat JWT token |
| `POST` | `/orders` | Bearer JWT | Kirim order ke orderbook |
| `GET` | `/orders` | Bearer JWT | Riwayat order lokal |
| `GET` | `/health` | — | Health check |

## Setup & Run

### 1. Struktur direktori yang diharapkan

```
project-root/
├── docker-compose.yml   ← file yang sudah diupdate (gabungan)
├── orderbook/           ← repo orderbook asli
└── trader/              ← folder mini project ini
```

### 2. Dapatkan API Key dari orderbook

Sebelum service trader bisa place order, kamu perlu buat API key di orderbook admin.

Buat client dulu lewat admin panel atau langsung via API:

```bash
# Login admin
curl -X POST http://localhost:3001/api/v1/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"yourpassword"}'

# Buat API key untuk client yang sudah ada
curl -X POST http://localhost:3001/api/v1/admin/api-keys \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{"client_id":"<client_uuid>","label":"trader-service","environment":"production"}'
```

Response akan berisi `key` dan `secret` — **secret hanya muncul sekali**.

### 3. Set environment variable

Edit `docker-compose.yml` di bagian service `trader`:

```yaml
ORDERBOOK_API_KEY: 'key-yang-kamu-dapat'
ORDERBOOK_API_SECRET: 'secret-yang-kamu-dapat'
```

### 4. Jalankan

```bash
docker compose up --build
```

### 5. Test API

```bash
# Register
curl -X POST http://localhost:9000/auth/register \
  -H "Content-Type: application/json" \
  -d '{"full_name":"Budi Trader","email":"budi@example.com","password":"secret123"}'

# Login
curl -X POST http://localhost:9000/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"budi@example.com","password":"secret123"}'
# → simpan token dari response

# Place Order
curl -X POST http://localhost:9000/orders \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTCUSDT","side":"bid","type":"limit","price":"90000","quantity":"0.01"}'

# List Orders
curl http://localhost:9000/orders \
  -H "Authorization: Bearer <token>"
```

## HMAC Authentication

Signature yang dikirim ke orderbook dihitung sesuai middleware orderbook:

```
payload  = METHOD + requestURI + timestamp + nonce
signature = HMAC-SHA256(payload, api_secret)
```

Header yang dikirim:
- `X-API-KEY`   — API key string
- `X-TIMESTAMP` — Unix timestamp (detik)
- `X-NONCE`     — UUID random per request
- `X-SIGNATURE` — hex(HMAC-SHA256(payload, secret))

## Database

Tabel dibuat otomatis di database `orderbook` yang sama (prefix `trader_`):

```sql
trader_users   -- user yang register ke trader service
trader_orders  -- riwayat order yang dikirim ke orderbook
```
