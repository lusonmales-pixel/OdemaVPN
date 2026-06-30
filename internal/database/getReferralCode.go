package database

import (
	"context"
	"math/rand/v2"
	"strconv"

	"github.com/jackc/pgx/v5"
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
			for i := 0; i < 6; i++ {
				num := rand.IntN(9)
				strNum := strconv.Itoa(num)
				FinalReferralCode += strNum
			}

			insertGenereatedCode := `
			INSERT INTO referral_codes (code, owner_id)
			VALUES ($1, $2)
			`
			_, err = conn.Exec(ctx, insertGenereatedCode, FinalReferralCode, tgID)
			if err != nil {
				return "", err
			}

			return FinalReferralCode, nil
		} else {
			return "", err
		}
	}

	return referralCode, nil
}
