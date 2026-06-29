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
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS subscriptions (
	user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
	status VARCHAR(20) DEFAULT 'inactive',
	expires_at TIMESTAMP DEFAULT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := conn.Exec(ctx, sqlQuery)
	if err != nil {
		return err
	}

	return nil
}
