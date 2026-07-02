package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func GetSubID(ctx context.Context, conn *pgx.Conn, TgID int64) (string, error) {
	var SubID string

	query := `
	SELECT sub_id FROM users
	WHERE id = $1
	`

	err := conn.QueryRow(ctx, query, TgID).Scan(&SubID)
	if err != nil {
		return "", err
	}

	return SubID, nil
}
