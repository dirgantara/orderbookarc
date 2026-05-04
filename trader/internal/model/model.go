package model

import "time"

type User struct {
	ID           string    `db:"id"            json:"id"`
	FullName     string    `db:"full_name"     json:"full_name"`
	Email        string    `db:"email"         json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	IsActive     bool      `db:"is_active"     json:"is_active"`
	CreatedAt    time.Time `db:"created_at"    json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"    json:"updated_at"`
}

type Order struct {
	ID          string    `db:"id"         json:"id"`
	UserID      string    `db:"user_id"    json:"user_id"`
	OrderCode   string    `db:"order_code"   json:"order_code"`
	Symbol      string    `db:"symbol"     json:"symbol"`
	Side        string    `db:"side"       json:"side"`
	Type        string    `db:"type"       json:"type"`
	Price       string    `db:"price"      json:"price"`
	Quantity    string    `db:"quantity"   json:"quantity"`
	Status      string    `db:"status"     json:"status"`
	RawResponse []byte    `db:"raw_response" json:"-"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
