package services

import (
	"avito-internship/internal/model"
	"avito-internship/internal/repository"
	"avito-internship/internal/repository/mocks"
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestSendCoin(t *testing.T) {
	tests := []struct {
		name            string
		fromUserID      int
		toUsername      string
		amount          int
		mockUser        func(mockUserRepo *mocks.MockUser)
		mockTransaction func(mockTransactionRepo *mocks.MockTransaction)
		expectedError   bool
	}{
		{
			name:       "successful transaction",
			fromUserID: 1,
			toUsername: "recipientUser",
			amount:     100,
			mockUser: func(mockUserRepo *mocks.MockUser) {
				mockUserRepo.EXPECT().
					GetUserByUsername(gomock.Any(), "recipientUser").
					Return(model.User{ID: 2}, nil).
					Times(1)
			},
			mockTransaction: func(mockTransactionRepo *mocks.MockTransaction) {
				mockTransactionRepo.EXPECT().
					SendCoin(gomock.Any(), 1, 2, 100).
					Return(nil).
					Times(1)
			},
			expectedError: false,
		},
		{
			name:       "user not found",
			fromUserID: 1,
			toUsername: "nonExistentUser",
			amount:     100,
			mockUser: func(mockUserRepo *mocks.MockUser) {
				mockUserRepo.EXPECT().
					GetUserByUsername(gomock.Any(), "nonExistentUser").
					Return(model.User{}, fmt.Errorf("user not found")).
					Times(1)
			},
			mockTransaction: func(mockTransactionRepo *mocks.MockTransaction) {},
			expectedError:   true,
		},
		{
			name:       "identical sender and recipient IDs",
			fromUserID: 1,
			toUsername: "senderUser",
			amount:     100,
			mockUser: func(mockUserRepo *mocks.MockUser) {
				mockUserRepo.EXPECT().
					GetUserByUsername(gomock.Any(), "senderUser").
					Return(model.User{ID: 1}, nil).
					Times(1)
			},
			mockTransaction: func(mockTransactionRepo *mocks.MockTransaction) {},
			expectedError:   true,
		},
		{
			name:       "transaction failure",
			fromUserID: 1,
			toUsername: "recipientUser",
			amount:     100,
			mockUser: func(mockUserRepo *mocks.MockUser) {
				mockUserRepo.EXPECT().
					GetUserByUsername(gomock.Any(), "recipientUser").
					Return(model.User{ID: 2}, nil).
					Times(1)
			},
			mockTransaction: func(mockTransactionRepo *mocks.MockTransaction) {
				mockTransactionRepo.EXPECT().
					SendCoin(gomock.Any(), 1, 2, 100).
					Return(fmt.Errorf("transaction error")).
					Times(1)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Мокируем репозитории
			mockTransactionRepo := mocks.NewMockTransaction(ctrl)
			mockUserRepo := mocks.NewMockUser(ctrl)

			// Запускаем моки для User и Transaction
			tt.mockUser(mockUserRepo)
			tt.mockTransaction(mockTransactionRepo)

			// Создаем сервис с моками
			repository := &repository.Repository{
				Transaction: mockTransactionRepo,
				User:        mockUserRepo,
				Buyer:       nil,
			}

			service := NewTransactionService(repository.Transaction, repository.User)

			// Вызываем метод SendCoin
			err := service.SendCoin(context.Background(), tt.fromUserID, tt.toUsername, tt.amount)

			// Проверяем ожидаемую ошибку
			if (err != nil) != tt.expectedError {
				t.Fatalf("expected error: %v, got: %v", tt.expectedError, err)
			}
		})
	}
}
