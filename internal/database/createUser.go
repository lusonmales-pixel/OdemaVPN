package database

import (
	"context"
	"math/rand/v2"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func CreateUserIfNotExits(ctx context.Context, conn *pgx.Conn, tgID int64, username string) (string, bool, error) {
	vless_uuid := uuid.New().String()
	var sub_id string
	allSymbolsStr := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321"
	for i := 0; i < 16; i++ {
		index := rand.IntN(len(allSymbolsStr))
		sub_id += string(allSymbolsStr[index])
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return "", false, err
	}

	defer tx.Rollback(ctx)

	query1 := `
	INSERT INTO users (id, username, vless_uuid, sub_id)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (id) DO NOTHING
	RETURNING id;
	`

	var InsertFlag int64

	err = tx.QueryRow(ctx, query1, tgID, username, vless_uuid, sub_id).Scan(&InsertFlag)

	isNew := true
	if err != nil {
		if err == pgx.ErrNoRows {
			isNew = false
		} else {
			return "", false, err
		}
	}

	if isNew {
		query2 := `
		INSERT INTO subscriptions (user_id, status, expires_at)
		VALUES ($1, 'active', NOW() + INTERVAL '1 day')
		ON CONFLICT (user_id) DO NOTHING
		`

		_, err = tx.Exec(ctx, query2, tgID)

		if err != nil {
			return "", false, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", false, err
	}

	return vless_uuid, true, nil

}
