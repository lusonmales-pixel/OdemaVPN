package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func InitTable(ctx context.Context, conn *pgx.Conn) error {
	sqlQuery := `
	CREATE TABLE IF NOT EXISTS users (
	id BIGINT PRIMARY KEY,
	username VARCHAR(100),
	vless_uuid VARCHAR(36) DEFAULT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	referrer_id BIGINT DEFAULT NULL,
	referral_bonus_given BOOL DEFAULT FALSE
	);

	CREATE TABLE IF NOT EXISTS subscriptions (
	user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
	status VARCHAR(20) DEFAULT 'inactive',
	expires_at TIMESTAMP DEFAULT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS referral_codes (
	code VARCHAR(6) PRIMARY KEY,
	owner_id BIGINT REFERENCES users(id) UNIQUE
	);
	`

	_, err := conn.Exec(ctx, sqlQuery)
	if err != nil {
		return err
	}

	return nil
}
