package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

func ApplyReffererBonus(ctx context.Context, conn *pgx.Conn, UserID int64) error {
	query1 := `
	SELECT referrer_id, referrer_bonus_given FROM users
	WHERE id = $1
	`

	query2 := `
	SELECT expires_at FROM subscriptions 
	WHERE user_id = $1
	`

	query3 := `
	UPDATE subscriptions
	SET expires_at = $1 + INTERVAL '7 days'
	WHERE user_id = $2
	`

	query4 := `
	UPDATE users
	SET referral_bonus_given = true
	WHERE id = $1
	`

	var RefferID *int64
	var BonusGiven bool
	var ExpiresAtNewUser time.Time
	var ExpiresAtReferrer time.Time

	err := conn.QueryRow(ctx, query1, UserID).Scan(&RefferID, &BonusGiven)
	if err != nil {
		return err
	}

	if RefferID == nil || BonusGiven {
		return nil
	}

	err = conn.QueryRow(ctx, query2, UserID).Scan(&ExpiresAtNewUser)
	if err != nil {
		return err
	}
	err = conn.QueryRow(ctx, query2, RefferID).Scan(&ExpiresAtReferrer)
	if err != nil {
		return err
	}

	_, err = conn.Exec(ctx, query3, ExpiresAtNewUser, UserID)
	if err != nil {
		return err
	}

	_, err = conn.Exec(ctx, query3, ExpiresAtReferrer, RefferID)
	if err != nil {
		return err
	}

	_, err = conn.Exec(ctx, query4, UserID)
	if err != nil {
		return err
	}

	return nil
}
