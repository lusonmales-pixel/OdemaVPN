package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func GetUserByRefCode(ctx context.Context, conn *pgx.Conn, refCode string, newUser int64) (int64, error) {
	getUser := `
	SELECT owner_id FROM referral_codes
	WHERE code = $1
	`
	var userId int64
	err := conn.QueryRow(ctx, getUser, refCode).Scan(&userId)
	if err != nil {
		return 0, err
	}

	addUserID := `
	UPDATE users
	SET referrer_id = $1
	WHERE id = $2
	`
	_, err = conn.Exec(ctx, addUserID, userId, newUser)
	if err != nil {
		return 0, err
	}

	return userId, nil
}
