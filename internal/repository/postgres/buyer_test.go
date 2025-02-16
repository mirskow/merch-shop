package postgres_test

import (
	"avito-internship/internal/model"
	"avito-internship/internal/repository/postgres"
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestBuy(t *testing.T) {
	// Инициализируем мок для базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании мока базы данных: %s", err)
	}
	defer db.Close()

	// Оборачиваем в sqlx.DB
	sqlxDB := sqlx.NewDb(db, "postgres")

	repo := postgres.NewBuyerPostgres(sqlxDB)

	type args struct {
		userID int
		item   model.Merch
	}

	tests := []struct {
		name    string
		args    args
		mock    func()
		wantErr bool
	}{
		{
			name: "Ok",
			args: args{
				userID: 1,
				item: model.Merch{
					ID:   1,
					Item: "t-shirt",
					Cost: 80,
				},
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"coin"}).AddRow(100)
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT coin FROM users WHERE id = \\$1 FOR UPDATE").WithArgs(1).WillReturnRows(rows)
				mock.ExpectExec("UPDATE users SET coin = coin - \\$1 WHERE id = \\$2").WithArgs(80, 1).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec("INSERT INTO purchases \\(user_id, item_id\\) VALUES \\(\\$1, \\$2\\)").WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Not Enough Coins",
			args: args{
				userID: 1,
				item: model.Merch{
					ID:   1,
					Item: "t-shirt",
					Cost: 200,
				},
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"coin"}).AddRow(100)
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT coin FROM users WHERE id = \\$1 FOR UPDATE").WithArgs(1).WillReturnRows(rows)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.Buy(context.Background(), tt.args.userID, tt.args.item)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
