package postgres_test

import (
	"avito-internship/internal/repository/postgres"
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSendCoin(t *testing.T) {
	// Инициализируем мок для базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании мока базы данных: %s", err)
	}
	defer db.Close()

	// Оборачиваем в sqlx.DB
	sqlxDB := sqlx.NewDb(db, "postgres")

	repo := postgres.NewTransactionPostgres(sqlxDB)

	type args struct {
		fromUser, toUser, amount int
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
				fromUser: 1,
				toUser:   2,
				amount:   50,
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"coin"}).AddRow(100)
				mock.ExpectBegin()
				// Экранируем символы $ в регулярных выражениях
				mock.ExpectQuery("SELECT coin FROM users WHERE id = \\$1 FOR UPDATE").WithArgs(1).WillReturnRows(rows)
				mock.ExpectExec("UPDATE users SET coin = coin - \\$1 WHERE id = \\$2").WithArgs(50, 1).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec("UPDATE users SET coin = coin \\+ \\$1 WHERE id = \\$2").WithArgs(50, 2).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec("INSERT INTO transactions \\(from_user, to_user, amount\\) VALUES \\(\\$1, \\$2, \\$3\\)").WithArgs(1, 2, 50).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Not Enough Coins",
			args: args{
				fromUser: 1,
				toUser:   2,
				amount:   200,
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"coin"}).AddRow(100)
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT coin FROM users WHERE id = \\$1 FOR UPDATE").WithArgs(1).WillReturnRows(rows)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Update Error",
			args: args{
				fromUser: 1,
				toUser:   2,
				amount:   50,
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"coin"}).AddRow(100)
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT coin FROM users WHERE id = \\$1 FOR UPDATE").WithArgs(1).WillReturnRows(rows)
				mock.ExpectExec("UPDATE users SET coin = coin - \\$1 WHERE id = \\$2").WithArgs(50, 1).WillReturnError(fmt.Errorf("update error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Insert Error",
			args: args{
				fromUser: 1,
				toUser:   2,
				amount:   50,
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"coin"}).AddRow(100)
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT coin FROM users WHERE id = \\$1 FOR UPDATE").WithArgs(1).WillReturnRows(rows)
				mock.ExpectExec("UPDATE users SET coin = coin - \\$1 WHERE id = \\$2").WithArgs(50, 1).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec("UPDATE users SET coin = coin \\+ \\$1 WHERE id = \\$2").WithArgs(50, 2).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec("INSERT INTO transactions \\(from_user, to_user, amount\\) VALUES \\(\\$1, \\$2, \\$3\\)").WithArgs(1, 2, 50).WillReturnError(fmt.Errorf("insert error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.SendCoin(context.Background(), tt.args.fromUser, tt.args.toUser, tt.args.amount)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
