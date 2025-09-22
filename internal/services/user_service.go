package services

import (
	"context"
	"time"

	"github.com/Quak1/gokei/internal/auth"
	"github.com/Quak1/gokei/internal/database/queries"
	"github.com/Quak1/gokei/internal/errors"
	"github.com/lib/pq"
)

type UserService struct {
	queries *queries.Queries
	jwt     *auth.JWTService
}

func NewUserService(queries *queries.Queries, jwtService *auth.JWTService) *UserService {
	return &UserService{
		queries: queries,
		jwt:     jwtService,
	}
}

type RegisterUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (s *UserService) Register(ctx context.Context, req *RegisterUserRequest) (*queries.User, error) {
	if req.Username == "" {
		return nil, errors.NewAppError(errors.ErrValidation, "Username is required", nil)
	}

	if req.Password == "" {
		return nil, errors.NewAppError(errors.ErrValidation, "Password is required", nil)
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, errors.NewAppError(errors.ErrInternal, "Failed to create user", err)
	}

	if req.Name == "" {
		req.Name = req.Username
	}

	user, err := s.queries.CreateUser(ctx, queries.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		Name:           req.Name,
	})
	if err != nil {
		if pqErr, _ := err.(*pq.Error); pqErr.Code == errors.PGErrorCodeUniqueViolation {
			return nil, errors.NewAppError(errors.ErrConflict, "Username is already in use", err)
		}
		return nil, errors.NewAppError(errors.ErrInternal, "Failed to create user", err)
	}

	return &user, nil
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *UserService) TokenLogin(ctx context.Context, req *LoginRequest) (string, error) {
	if req.Username == "" {
		return "", errors.NewAppError(errors.ErrValidation, "Username is required", nil)
	}
	if req.Password == "" {
		return "", errors.NewAppError(errors.ErrValidation, "Password is required", nil)
	}

	loginFailedMessage := "Login failed. Please check your credentials."

	user, err := s.queries.GetUser(ctx, req.Username)
	if err != nil {
		return "", errors.NewAppError(errors.ErrInternal, loginFailedMessage, err)
	}

	err = auth.CheckHashedPassword(req.Password, user.HashedPassword)
	if err != nil {
		return "", errors.NewAppError(errors.ErrInternal, loginFailedMessage, err)
	}

	token, err := s.jwt.MakeJWT(user, time.Hour*24)
	if err != nil {
		return "", errors.NewAppError(errors.ErrInternal, loginFailedMessage, err)
	}

	return token, nil
}
