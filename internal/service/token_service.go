package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"

	"github.com/Quak1/gokei/internal/database"
	"github.com/Quak1/gokei/internal/database/store"
	"github.com/Quak1/gokei/pkg/validator"
)

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	Expiry    time.Time `json:"expiry"`
}

type TokenService struct {
	queries *store.Queries
}

func NewTokenService(queries *store.Queries) *TokenService {
	return &TokenService{
		queries: queries,
	}
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

func HashToken(token string) []byte {
	hash := sha256.Sum256([]byte(token))
	return hash[:]
}

func (s *TokenService) New(userID int, ttl time.Duration) (*Token, error) {
	token := &Token{
		Expiry: time.Now().Add(ttl),
	}

	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(bytes)
	token.Hash = HashToken(token.Plaintext)

	_, err = s.queries.CreateToken(context.Background(), store.CreateTokenParams{
		Hash:   token.Hash,
		UserID: int32(userID),
		Expiry: token.Expiry,
	})
	if err != nil {
		return nil, database.HandleForeignKeyError(err)
	}

	return token, nil
}
