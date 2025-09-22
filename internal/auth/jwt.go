package auth

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Quak1/gokei/internal/database/queries"
	"github.com/Quak1/gokei/internal/errors"
	"github.com/Quak1/gokei/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

const tokenIssuer = "gokei-access"

type CustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type JWTService struct {
	tokenSecret []byte
}

func NewJWTService(tokenSecret []byte) *JWTService {
	return &JWTService{
		tokenSecret: tokenSecret,
	}
}

func (s *JWTService) MakeJWT(user queries.User, expiresIn time.Duration) (string, error) {
	claims := CustomClaims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    tokenIssuer,
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
			Subject:   strconv.Itoa(int(user.ID)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.tokenSecret)
}

func (s *JWTService) validateJWT(tokenString string) (*CustomClaims, error) {
	errorMessage := "Invalid token"

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (any, error) {
		return s.tokenSecret, nil
	}, jwt.WithIssuer(tokenIssuer))
	if err != nil {
		return nil, errors.NewAppError(errors.ErrUnauthorized, errorMessage, err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.NewAppError(errors.ErrUnauthorized, errorMessage, err)
	}

	return claims, nil
}

func (s *JWTService) getBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.NewAppError(errors.ErrUnauthorized, "Missing Authentication header", nil)
	}

	split := strings.Split(authHeader, " ")
	if len(split) != 2 || split[0] != "Bearer" {
		return "", errors.NewAppError(errors.ErrUnauthorized, "Invalid Authorization header", nil)
	}

	return split[1], nil
}

func (s *JWTService) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := s.getBearerToken(r.Header)
		if err != nil {
			utils.ResError(w, err)
			return
		}

		claims, err := s.validateJWT(token)
		if err != nil {
			utils.ResError(w, err)
			return
		}

		newCtx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}
