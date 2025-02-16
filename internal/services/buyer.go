package services

import (
	"avito-internship/internal/repository"
	"context"
)

type BuyerService struct {
	repo repository.Buyer
}

func NewBuyerService(repo repository.Buyer) *BuyerService {
	return &BuyerService{
		repo: repo,
	}
}

func (b *BuyerService) Buy(ctx context.Context, userID int, itemName string) error {
	item, err := b.repo.GetItem(ctx, itemName)
	if err != nil {
		return err
	}
	return b.repo.Buy(ctx, userID, item)
}
