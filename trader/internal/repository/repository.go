package repository

import (
	"context"
	"trader/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *model.User) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO trader_users (id, full_name, email, password_hash, is_active, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, NOW(), NOW())`,
		u.ID, u.FullName, u.Email, u.PasswordHash, u.IsActive,
	)
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	u := &model.User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, full_name, email, password_hash, is_active, created_at, updated_at
		 FROM trader_users WHERE email = $1`, email,
	).Scan(&u.ID, &u.FullName, &u.Email, &u.PasswordHash, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// ─── Order Repository ──────────────────────────────────────────────────────────

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, o *model.Order) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO trader_orders (id, user_id, order_code, symbol, side, type, price, quantity, status, raw_response, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW())`,
		o.ID, o.UserID, o.OrderCode, o.Symbol, o.Side, o.Type, o.Price, o.Quantity, o.Status, o.RawResponse,
	)
	return err
}

func (r *OrderRepository) FindByUserID(ctx context.Context, userID string) ([]model.Order, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, order_code, symbol, side, type, price, quantity, status, created_at
		 FROM trader_orders WHERE user_id = $1 ORDER BY created_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.OrderCode, &o.Symbol, &o.Side, &o.Type,
			&o.Price, &o.Quantity, &o.Status, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *OrderRepository) FindByOrderCode(ctx context.Context, userID, orderCode string) (*model.Order, error) {
	o := &model.Order{}
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, order_code, symbol, side, type, price, quantity, status, created_at
		 FROM trader_orders WHERE user_id = $1 AND order_code = $2`,
		userID, orderCode,
	).Scan(
		&o.ID, &o.UserID, &o.OrderCode, &o.Symbol, &o.Side, &o.Type,
		&o.Price, &o.Quantity, &o.Status, &o.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return o, nil
}

// UpdateStatusByOrderCode — dipanggil saat webhook dari orderbook diterima.
func (r *OrderRepository) UpdateStatusByOrderCode(ctx context.Context, orderCode, status string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE trader_orders SET status = $1 WHERE order_code = $2`,
		status, orderCode,
	)
	return err
}
