package services

import (
	"avito-internship/internal/repository"
	"context"
	"fmt"
)

type TransactionService struct {
	repo     repository.Transaction
	userRepo repository.User
}

func NewTransactionService(repo repository.Transaction, userRepo repository.User) *TransactionService {
	return &TransactionService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (ts *TransactionService) SendCoin(ctx context.Context, fromUserID int, toUsername string, amount int) error {
	toUser, err := ts.userRepo.GetUserByUsername(ctx, toUsername)
	if err != nil {
		return err
	}

	if fromUserID == toUser.ID {
		return fmt.Errorf("identical sender and recipient IDs")
	}

	return ts.repo.SendCoin(ctx, fromUserID, toUser.ID, amount)
}
