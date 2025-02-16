package services

import (
	"avito-internship/internal/model"
	"avito-internship/internal/repository"
	"avito-internship/internal/repository/mocks"
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestGetUsersInfo(t *testing.T) {
	tests := []struct {
		name          string
		userID        int
		mockUserRepo  func(mockUserRepo *mocks.MockUser)
		expectedCoins int
		expectedPurch []model.Purchases
		expectedSent  []model.Transaction
		expectedRecv  []model.Transaction
		expectedError bool
	}{
		{
			name:   "successful retrieval with transactions",
			userID: 1,
			mockUserRepo: func(mockUserRepo *mocks.MockUser) {
				mockUserRepo.EXPECT().
					GetUsersInfo(gomock.Any(), 1).
					Return(1000, []model.Purchases{
						{ID: 1, ItemID: 10},
					}, []model.Transaction{
						{FromUserID: 1, ToUserID: 2, Amount: 50},
						{FromUserID: 3, ToUserID: 1, Amount: 30},
					}, nil).
					Times(1)
			},
			expectedCoins: 1000,
			expectedPurch: []model.Purchases{
				{ID: 1, ItemID: 10},
			},
			expectedSent: []model.Transaction{
				{FromUserID: 1, ToUserID: 2, Amount: 50},
			},
			expectedRecv: []model.Transaction{
				{FromUserID: 3, ToUserID: 1, Amount: 30},
			},
			expectedError: false,
		},
		{
			name:   "user without transactions",
			userID: 3,
			mockUserRepo: func(mockUserRepo *mocks.MockUser) {
				mockUserRepo.EXPECT().
					GetUsersInfo(gomock.Any(), 3).
					Return(500, []model.Purchases{}, []model.Transaction{}, nil).
					Times(1)
			},
			expectedCoins: 500,
			expectedPurch: []model.Purchases{},
			expectedSent:  []model.Transaction{},
			expectedRecv:  []model.Transaction{},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mocks.NewMockUser(ctrl)
			tt.mockUserRepo(mockUserRepo)

			repository := &repository.Repository{
				User: mockUserRepo,
			}

			service := NewUserService(repository.User)

			coins, purchases, sent, recv, err := service.GetUsersInfo(context.Background(), tt.userID)

			if (err != nil) != tt.expectedError {
				t.Fatalf("expected error: %v, got: %v", tt.expectedError, err)
			}

			if coins != tt.expectedCoins {
				t.Errorf("expected coins: %d, got: %d", tt.expectedCoins, coins)
			}

			// Нормализуем nil в пустые слайсы перед сравнением
			if sent == nil {
				sent = []model.Transaction{}
			}
			if recv == nil {
				recv = []model.Transaction{}
			}
			if purchases == nil {
				purchases = []model.Purchases{}
			}

			if !reflect.DeepEqual(purchases, tt.expectedPurch) {
				t.Errorf("expected purchases: %+v, got: %+v", tt.expectedPurch, purchases)
			}

			if !reflect.DeepEqual(sent, tt.expectedSent) {
				t.Errorf("expected sent transactions: %+v, got: %+v", tt.expectedSent, sent)
			}

			if !reflect.DeepEqual(recv, tt.expectedRecv) {
				t.Errorf("expected received transactions: %+v, got: %+v", tt.expectedRecv, recv)
			}
		})
	}
}
