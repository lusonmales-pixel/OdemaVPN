package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func GetUUID(ctx context.Context, conn *pgx.Conn, tgID int64) (uuid string, err error) {
	query := `
	SELECT vless_uuid FROM users
	WHERE id = $1	
	`

	rows, err := conn.Query(ctx, query, tgID)
	if err != nil {
		return "", err
	}

	for rows.Next() {
		err = rows.Scan(&uuid)
		if err != nil {
			return "", err
		}
	}

	return uuid, nil

}
