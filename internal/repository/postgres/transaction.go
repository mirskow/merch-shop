package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type TransactionPostgres struct {
	db *sqlx.DB
}

func NewTransactionPostgres(db *sqlx.DB) *TransactionPostgres {
	return &TransactionPostgres{db: db}
}

func (tp *TransactionPostgres) SendCoin(ctx context.Context, fromUser, toUser, amount int) error {
	tx, err := tp.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		var e error
		if err == nil {
			e = tx.Commit()
		} else {
			e = tx.Rollback()
		}

		if err == nil && e != nil {
			err = fmt.Errorf("finishing transaction: %w", e)
		}
	}()

	var balance int
	err = tx.QueryRowContext(ctx, "SELECT coin FROM users WHERE id = $1 FOR UPDATE", fromUser).Scan(&balance)
	if err != nil {
		return err
	}

	if balance < amount {
		err = fmt.Errorf("not enough coins")
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE users SET coin = coin - $1 WHERE id = $2", amount, fromUser)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE users SET coin = coin + $1 WHERE id = $2", amount, toUser)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO transactions (from_user, to_user, amount) VALUES ($1, $2, $3)", fromUser, toUser, amount)
	if err != nil {
		return err
	}

	return nil
}
