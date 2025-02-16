package services

import (
	"avito-internship/internal/model"
	"avito-internship/internal/repository/mocks"
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestBuy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Мокаем репозиторий
	mockBuyerRepo := mocks.NewMockBuyer(ctrl)

	// Определяем тестовые сценарии
	tests := []struct {
		name        string
		userID      int
		itemName    string
		mockSetup   func()
		expectError bool
	}{
		{
			name:     "successful purchase",
			userID:   1,
			itemName: "t-shirt",
			mockSetup: func() {
				mockBuyerRepo.EXPECT().
					GetItem(gomock.Any(), "t-shirt").
					Return(model.Merch{ID: 1, Item: "t-shirt", Cost: 80}, nil).Times(1)

				mockBuyerRepo.EXPECT().
					Buy(gomock.Any(), 1, model.Merch{ID: 1, Item: "t-shirt", Cost: 80}).
					Return(nil).Times(1)
			},
			expectError: false,
		},
		{
			name:     "item not found",
			userID:   1,
			itemName: "non-existent-item",
			mockSetup: func() {
				mockBuyerRepo.EXPECT().
					GetItem(gomock.Any(), "non-existent-item").
					Return(model.Merch{}, errors.New("item not found")).Times(1)
			},
			expectError: true,
		},
		{
			name:     "purchase fails",
			userID:   1,
			itemName: "cup",
			mockSetup: func() {
				mockBuyerRepo.EXPECT().
					GetItem(gomock.Any(), "cup").
					Return(model.Merch{ID: 2, Item: "cup", Cost: 20}, nil).Times(1)

				mockBuyerRepo.EXPECT().
					Buy(gomock.Any(), 1, model.Merch{ID: 2, Item: "cup", Cost: 20}).
					Return(errors.New("purchase failed")).Times(1)
			},
			expectError: true,
		},
	}

	// Запускаем тесты
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			service := NewBuyerService(mockBuyerRepo)
			err := service.Buy(context.Background(), tc.userID, tc.itemName)

			if (err != nil) != tc.expectError {
				t.Errorf("expected error: %v, got: %v", tc.expectError, err)
			}
		})
	}
}
