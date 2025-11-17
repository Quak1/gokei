package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Quak1/gokei/internal/database"
	"github.com/Quak1/gokei/internal/database/store"
	"github.com/Quak1/gokei/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid authentication credentials")
)

func doesPasswordMatch(user *store.User, plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

type AuthService struct {
	queries      store.QuerierTx
	tokenService *TokenService
}

func NewAuthService(queries store.QuerierTx, tokenService *TokenService) *AuthService {
	return &AuthService{
		queries:      queries,
		tokenService: tokenService,
	}
}

func (s *AuthService) CreateAuthToken(username, password string) (*Token, error) {
	v := validator.New()
	validateUsername(v, username)
	validatePasswordPlaintext(v, password)
	if !v.Valid() {
		return nil, v.GetErrors()
	}

	user, err := s.queries.GetUserByUsername(context.Background(), username)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, database.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	match, err := doesPasswordMatch(&user, password)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, ErrInvalidCredentials
	}

	token, err := s.tokenService.New(int(user.ID), 24*time.Hour)
	if err != nil {
		return nil, err
	}

	return token, nil
}
