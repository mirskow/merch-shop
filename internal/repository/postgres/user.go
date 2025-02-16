package postgres

import (
	"avito-internship/internal/model"
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) CreateUser(ctx context.Context, user model.User) (int, error) {
	var id int

	query := "INSERT INTO users (username, password_hash, coin) VALUES ($1, $2, DEFAULT) RETURNING id"
	err := r.db.QueryRowContext(ctx, query, user.Username, user.PasswordHash).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (u *UserPostgres) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	return u.getUserByField(ctx, "username", username)
}

func (u *UserPostgres) GetUserByID(ctx context.Context, userID int) (model.User, error) {
	return u.getUserByField(ctx, "id", userID)
}

func (u *UserPostgres) getUserByField(ctx context.Context, field string, value any) (model.User, error) {
	var user model.User

	query := fmt.Sprintf("SELECT * FROM users WHERE %s=$1", field)
	err := u.db.GetContext(ctx, &user, query, value)

	return user, err
}

func (u *UserPostgres) GetUsersInfo(ctx context.Context, userID int) (int, []model.Purchases, []model.Transaction, error) {
	tx, err := u.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return 0, nil, nil, err
	}
	defer tx.Rollback()

	coins, err := u.getUsersCoins(ctx, tx, userID)
	if err != nil {
		return 0, nil, nil, err
	}

	purchases, err := u.getUsersPurchases(ctx, tx, userID)
	if err != nil {
		return 0, nil, nil, err
	}

	historyTransaction, err := u.getTransactionHistory(ctx, tx, userID)
	if err != nil {
		return 0, nil, nil, err
	}

	if err := tx.Commit(); err != nil {
		return 0, nil, nil, err
	}

	return coins, purchases, historyTransaction, nil
}

func (u *UserPostgres) getUsersCoins(ctx context.Context, tx *sql.Tx, userID int) (int, error) {
	var coins int
	err := tx.QueryRowContext(ctx, "SELECT coin FROM users WHERE id = $1", userID).Scan(&coins)
	if err != nil {
		logrus.Error("Error fetching coins:", err)
		return 0, err
	}
	return coins, nil
}

func (u *UserPostgres) getUsersPurchases(ctx context.Context, tx *sql.Tx, userID int) ([]model.Purchases, error) {
	var purchases []model.Purchases
	query := "SELECT t.item as item, COUNT(p.id) as quantity FROM purchases p JOIN merch t ON p.item_id = t.id WHERE p.user_id = $1 GROUP BY t.item"

	rows, err := tx.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.Purchases

		if err := rows.Scan(&item.ItemName, &item.Quantity); err != nil {
			return nil, err
		}

		purchases = append(purchases, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return purchases, nil
}

func (u *UserPostgres) getTransactionHistory(ctx context.Context, tx *sql.Tx, userID int) ([]model.Transaction, error) {
	transactions := make([]model.Transaction, 0)

	query := `SELECT 
			t.from_user, t.to_user, t.amount, 
			sender.username, receiver.username
		FROM transactions t
		JOIN users sender ON t.from_user = sender.id
		JOIN users receiver ON t.to_user = receiver.id
		WHERE t.from_user = $1 OR t.to_user = $2`

	rows, err := tx.QueryContext(ctx, query, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction model.Transaction
		if err := rows.Scan(&transaction.FromUserID, &transaction.ToUserID, &transaction.Amount, &transaction.SenderName, &transaction.ReceiverName); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
