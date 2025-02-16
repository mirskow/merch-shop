package services

import (
	"avito-internship/internal/model"
	"avito-internship/internal/repository"
	"context"
)

type Authorization interface {
	CreateUser(ctx context.Context, user model.User) (string, error)
	GenerateToken(id int) (string, error)
	ParseToken(token string) (int, error)
}

type Transaction interface {
	SendCoin(ctx context.Context, fromUser int, toUser string, amount int) error
}

type User interface {
	GetUsersInfo(ctx context.Context, userID int) (int, []model.Purchases, []model.Transaction, []model.Transaction, error)
}

type Buyer interface {
	Buy(ctx context.Context, userID int, thing string) error
}

type Services struct {
	Authorization
	Transaction
	User
	Buyer
}

func NewServices(repos *repository.Repository) *Services {
	return &Services{
		Authorization: NewAuthService(repos.User),
		Transaction:   NewTransactionService(repos.Transaction, repos.User),
		User:          NewUserService(repos.User),
		Buyer:         NewBuyerService(repos.Buyer),
	}
}
