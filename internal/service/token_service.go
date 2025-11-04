package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"

	"github.com/Quak1/gokei/internal/database"
	"github.com/Quak1/gokei/internal/database/store"
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
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

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
