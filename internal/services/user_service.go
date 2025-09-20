package services

import (
	"context"

	"github.com/Quak1/gokei/internal/auth"
	"github.com/Quak1/gokei/internal/database/queries"
	"github.com/Quak1/gokei/internal/errors"
	"github.com/lib/pq"
)

type UserService struct {
	queries *queries.Queries
}

func NewUserService(queries *queries.Queries) *UserService {
	return &UserService{
		queries: queries,
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
