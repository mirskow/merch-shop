package services

import (
	"avito-internship/internal/model"
	"avito-internship/internal/repository"
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	tokenTTL = 12 * time.Hour
)

type AuthService struct {
	repo repository.User
}

type tokenClaims struct {
	jwt.StandardClaims
	UserID int `json:"user_id"`
}

func NewAuthService(repo repository.User) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(ctx context.Context, user model.User) (string, error) {
	user.PasswordHash = s.generatePasswordHash(user.PasswordHash)

	userFromDb, err := s.repo.GetUserByUsername(ctx, user.Username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("user check error: %w", err)
	}

	if errors.Is(err, sql.ErrNoRows) {
		user.ID, err = s.repo.CreateUser(ctx, user)
		if err != nil {
			return "", fmt.Errorf("user creation error: %w", err)
		}
	} else {
		user.ID = userFromDb.ID
	}

	token, err := s.GenerateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("token generation error: %w", err)
	}

	return token, nil
}

func (s *AuthService) ParseToken(accesstoken string) (int, error) {
	token, err := jwt.ParseWithClaims(accesstoken, &tokenClaims{}, func(token *jwt.Token) (any, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims do not belong to type *tokenClaims")
	}

	return claims.UserID, nil
}

func (s *AuthService) GenerateToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&tokenClaims{
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(tokenTTL).Unix(),
				IssuedAt:  time.Now().Unix(),
			},
			userId,
		})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(os.Getenv("PASSWORD_SALT"))))
}
