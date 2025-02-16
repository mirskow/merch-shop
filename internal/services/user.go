package services

import (
	"avito-internship/internal/model"
	"avito-internship/internal/repository"
	"context"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (u *UserService) GetUsersInfo(ctx context.Context, userID int) (int, []model.Purchases, []model.Transaction, []model.Transaction, error) {
	coins, purchases, historyTransaction, err := u.repo.GetUsersInfo(ctx, userID)
	if err != nil {
		return 0, nil, nil, nil, err
	}

	var sent []model.Transaction
	var receiver []model.Transaction

	for _, t := range historyTransaction {
		if t.FromUserID == userID {
			sent = append(sent, t)
		} else {
			receiver = append(receiver, t)
		}
	}

	return coins, purchases, sent, receiver, nil
}
