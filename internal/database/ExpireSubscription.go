package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type ExpiredUser struct {
	Vless_uuid string
	TgId       int64
}

func Expire(ctx context.Context, conn *pgx.Conn) ([]ExpiredUser, error) {
	query := `
	UPDATE subscriptions s
	SET status = 'inactive'
	FROM users u
	WHERE s.user_id = u.id 
	  AND s.status = 'active' 
	  AND s.expires_at < NOW()
	RETURNING u.vless_uuid, u.id;
	`

	rows, err := conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var expiredUsers []ExpiredUser
	for rows.Next() {
		var eu ExpiredUser
		err := rows.Scan(&eu.Vless_uuid, &eu.TgId)
		if err != nil {
			return nil, err
		}
		expiredUsers = append(expiredUsers, eu)
	}

	return expiredUsers, nil
}
