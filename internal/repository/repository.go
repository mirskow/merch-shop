package repository

import (
	"avito-internship/internal/model"
	"avito-internship/internal/repository/postgres"
	"context"

	"github.com/jmoiron/sqlx"
)

type Transaction interface {
	SendCoin(ctx context.Context, fromUser int, toUser int, amount int) error
}

type User interface {
	CreateUser(ctx context.Context, user model.User) (int, error)
	GetUsersInfo(ctx context.Context, userID int) (int, []model.Purchases, []model.Transaction, error)
	GetUserByUsername(ctx context.Context, username string) (model.User, error)
	GetUserByID(ctx context.Context, userID int) (model.User, error)
}

type Buyer interface {
	Buy(ctx context.Context, userID int, item model.Merch) error
	GetItem(ctx context.Context, itemName string) (model.Merch, error)
}

type Repository struct {
	Transaction
	User
	Buyer
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Transaction: postgres.NewTransactionPostgres(db),
		User:        postgres.NewUserPostgres(db),
		Buyer:       postgres.NewBuyerPostgres(db),
	}
}
