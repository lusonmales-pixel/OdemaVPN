package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func CheckStatus(ctx context.Context, conn *pgx.Conn, userID int64) (status string, err error) {

	query := `
	SELECT status FROM subscriptions
	WHERE user_id = $1
	`

	err = conn.QueryRow(ctx, query, userID).Scan(&status)
	if err != nil {
		return "", err
	}

	return status, nil
}
