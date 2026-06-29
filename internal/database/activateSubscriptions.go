package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

func DefineMounthsAmount(amount float64) int {
	if amount == 150.0 {
		return 1
	}
	if amount == 400.0 {
		return 3
	}
	if amount == 1300 {
		return 12
	}

	return 0
}

func Activate(ctx context.Context, conn *pgx.Conn, tgID int64, amount float64) error {

	mounthAmount := DefineMounthsAmount(amount)
	if mounthAmount == 0 {
		log.Println("Wrong payment!")
		err := errors.New("Wrong payment")
		return err
	}

	var currentExpireDate time.Time
	var newExpireDate time.Time

	getExpireDate := `
	SELECT expires_at
	FROM subscriptions
	WHERE user_id = $1
	`

	activateSubscription := `
	UPDATE subscriptions 
	SET expires_at = $1, status = 'active'
	WHERE user_id = $2
	`

	err := conn.QueryRow(ctx, getExpireDate, tgID).Scan(&currentExpireDate)
	if err != nil {
		log.Println("Error", err)
		return err
	}

	if currentExpireDate.After(time.Now()) {
		newExpireDate = currentExpireDate.AddDate(0, mounthAmount, 0)
	} else {
		newExpireDate = time.Now().AddDate(0, mounthAmount, 0)
	}

	_, err = conn.Exec(ctx, activateSubscription, newExpireDate, tgID)
	if err != nil {
		log.Println("Error:", err)
	}

	return nil

}
