package postgres_test

import (
	"avito-internship/internal/model"
	"avito-internship/internal/repository/postgres"
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	// Инициализируем мок для базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании мока базы данных: %s", err)
	}
	defer db.Close()

	// Оборачиваем в sqlx.DB
	sqlxDB := sqlx.NewDb(db, "postgres")

	repo := postgres.NewUserPostgres(sqlxDB)

	tests := []struct {
		name    string
		user    model.User
		mock    func()
		wantID  int
		wantErr bool
	}{
		{
			name: "Ok",
			user: model.User{
				Username:     "testuser",
				PasswordHash: "hashedpassword",
			},
			mock: func() {
				// Ожидаем запрос на вставку пользователя
				mock.ExpectQuery("INSERT INTO users \\(username, password_hash, coin\\) VALUES \\(\\$1, \\$2, DEFAULT\\) RETURNING id").
					WithArgs("testuser", "hashedpassword").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantID:  1,
			wantErr: false,
		},
		{
			name: "Error inserting user",
			user: model.User{
				Username:     "testuser",
				PasswordHash: "hashedpassword",
			},
			mock: func() {
				// Ожидаем ошибку при вставке
				mock.ExpectQuery("INSERT INTO users \\(username, password_hash, coin\\) VALUES \\(\\$1, \\$2, DEFAULT\\) RETURNING id").
					WithArgs("testuser", "hashedpassword").
					WillReturnError(fmt.Errorf("insert error"))
			},
			wantID:  0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			id, err := repo.CreateUser(context.Background(), tt.user)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantID, id)
			}
			// Проверяем, что все ожидания были выполнены
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
