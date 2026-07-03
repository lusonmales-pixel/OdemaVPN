package database

import (
	"context"
	"errors"
	"math/rand/v2"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func GetReferralCode(ctx context.Context, conn *pgx.Conn, tgID int64) (string, error) {
	getCode := `
	SELECT code FROM referral_codes
	WHERE owner_id = $1
	`
	var FinalReferralCode string
	var referralCode string
	err := conn.QueryRow(ctx, getCode, tgID).Scan(&referralCode)
	if err != nil {
		if err == pgx.ErrNoRows {
			insertGenereatedCode := `
			INSERT INTO referral_codes (code, owner_id)
			VALUES ($1, $2)
			`

			for attempt := 0; attempt < 10; attempt++ {
				FinalReferralCode = ""
				for i := 0; i < 6; i++ {
					num := rand.IntN(10)
					strNum := strconv.Itoa(num)
					FinalReferralCode += strNum
				}

				_, err = conn.Exec(ctx, insertGenereatedCode, FinalReferralCode, tgID)
				if err == nil {
					return FinalReferralCode, nil
				}

				var pgErr *pgconn.PgError
				if errors.As(err, &pgErr) && pgErr.Code == "23505" {
					continue
				}

				return "", err
			}

			return "", pgx.ErrNoRows
		} else {
			return "", err
		}
	}

	return referralCode, nil
}
