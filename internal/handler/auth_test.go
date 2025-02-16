package handler

import (
	"avito-internship/internal/model"
	"avito-internship/internal/services"
	service_mocks "avito-internship/internal/services/mocks"
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_auth(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *service_mocks.MockAuthorization, user model.User)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            model.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"username": "username", "password": "qwerty"}`,
			inputUser: model.User{
				Username:     "username",
				PasswordHash: "qwerty",
			},
			mockBehavior: func(r *service_mocks.MockAuthorization, user model.User) {
				r.EXPECT().CreateUser(gomock.Any(), user).Return("mockToken", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"token":"mockToken"}`,
		},
		{
			name:                 "Wrong Input",
			inputBody:            `{"username": "username"}`,
			inputUser:            model.User{}, // Пустой user
			mockBehavior:         func(r *service_mocks.MockAuthorization, user model.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"Неверный запрос."}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"username": "username", "password": "qwerty"}`,
			inputUser: model.User{
				Username:     "username",
				PasswordHash: "qwerty",
			},
			mockBehavior: func(r *service_mocks.MockAuthorization, user model.User) {
				// Ошибка при вызове CreateUser
				r.EXPECT().CreateUser(gomock.Any(), user).Return("", errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"Внутренняя ошибка сервера."}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := service_mocks.NewMockAuthorization(c)
			test.mockBehavior(repo, test.inputUser)

			services := &services.Services{Authorization: repo}
			handler := Handler{services}

			// Init Endpoint
			r := gin.New()
			r.POST("/auth", handler.auth)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth", bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
