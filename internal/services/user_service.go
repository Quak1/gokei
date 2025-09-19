package services

import (
	"context"
	"fmt"

	"github.com/Quak1/gokei/internal/auth"
	"github.com/Quak1/gokei/internal/database/queries"
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
		return nil, fmt.Errorf("username is required")
	}

	if req.Password == "" {
		return nil, fmt.Errorf("password is required")
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
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
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}
