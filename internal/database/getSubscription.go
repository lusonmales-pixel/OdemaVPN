package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type SubscriptionData struct {
	Status    string
	ExpiresAt *time.Time
}

func GetSubscription(ctx context.Context, conn *pgx.Conn, userID int64) (SubscriptionData, error) {
	query := `
	SELECT status, expires_at
	FROM subscriptions
	WHERE user_id = $1
	`

	var subscription SubscriptionData
	err := conn.QueryRow(ctx, query, userID).Scan(&subscription.Status, &subscription.ExpiresAt)
	if err != nil {
		return SubscriptionData{}, err
	}

	return subscription, nil
}
