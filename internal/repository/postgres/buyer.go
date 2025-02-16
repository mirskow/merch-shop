package postgres

import (
	"avito-internship/internal/model"
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type BuyerPostgres struct {
	db *sqlx.DB
}

func NewBuyerPostgres(db *sqlx.DB) *BuyerPostgres {
	return &BuyerPostgres{db: db}
}

func (b *BuyerPostgres) Buy(ctx context.Context, userID int, item model.Merch) error {
	tx, err := b.db.BeginTx(ctx, &sql.TxOptions{})
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
	err = tx.QueryRowContext(ctx, "SELECT coin FROM users WHERE id = $1 FOR UPDATE", userID).Scan(&balance)
	if err != nil {
		return err
	}

	if balance < item.Cost {
		err = fmt.Errorf("not enough coins")
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE users SET coin = coin - $1 WHERE id = $2", item.Cost, userID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO purchases (user_id, item_id) VALUES ($1, $2)", userID, item.ID)
	if err != nil {
		return err
	}

	return nil
}

func (b *BuyerPostgres) GetItem(ctx context.Context, itemName string) (model.Merch, error) {
	var item model.Merch

	err := b.db.GetContext(ctx, &item, "SELECT * FROM merch WHERE item=$1", itemName)

	return item, err
}
