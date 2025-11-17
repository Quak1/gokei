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
	ErrDuplicateUsername = errors.New("duplicate username")
)

func validatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(validator.NonZero(password), "password", "Must be provided")
	v.Check(validator.MinLength(password, 8), "password", "Must be at least 8 bytes long")
	v.Check(validator.MaxLength(password, 72), "password", "Must not be more than 72 bytes long")
}

func validateName(v *validator.Validator, username string) {
	v.Check(validator.NonZero(username), "name", "Must be provided")
	v.Check(validator.MinLength(username, 2), "name", "Must be at least 2 bytes long")
	v.Check(validator.MaxLength(username, 100), "name", "Must not be more than 100 bytes long")
}

func validateUsername(v *validator.Validator, username string) {
	v.Check(validator.NonZero(username), "username", "Must be provided")
	v.Check(validator.MinLength(username, 2), "username", "Must be at least 2 bytes long")
	v.Check(validator.MaxLength(username, 20), "username", "Must not be more than 20 bytes long")
}

func validateUser(v *validator.Validator, user *InputUser) {
	validateUsername(v, user.Username)
	validateName(v, user.Name)
	validatePasswordPlaintext(v, user.Password)
}

func hashPassword(plaintextPassword string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

type UserService struct {
	queries store.QuerierTx
}

func NewUserService(queries store.QuerierTx) *UserService {
	return &UserService{
		queries: queries,
	}
}

type InputUser struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (s *UserService) Create(params *InputUser) (*store.User, error) {
	v := validator.New()
	if validateUser(v, params); !v.Valid() {
		return nil, v.GetErrors()
	}

	hash, err := hashPassword(params.Password)
	if err != nil {
		return nil, err
	}

	user := &store.CreateUserParams{
		Username:     params.Username,
		Name:         params.Name,
		PasswordHash: hash,
	}

	data, err := s.queries.CreateUser(context.Background(), *user)
	if err != nil {
		if database.IsUniqueContraintViolation(err) {
			return nil, ErrDuplicateUsername
		}
		return nil, err
	}

	return &data, nil
}

func (s *UserService) GetByID(id int32) (*store.User, error) {
	if id < 1 {
		return nil, database.ErrRecordNotFound
	}

	user, err := s.queries.GetUserByID(context.Background(), id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, database.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (s *UserService) DeleteByID(id int32) error {
	if id < 1 {
		return database.ErrRecordNotFound
	}

	result, err := s.queries.DeleteUserById(context.Background(), id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return database.ErrRecordNotFound
	}

	return nil
}

type UpdateUserParams struct {
	Name     *string `json:"name"`
	Password *string `json:"password"`
}

func (s *UserService) UpdateByID(id int32, updateParams *UpdateUserParams) (*store.User, error) {
	if id < 1 {
		return nil, database.ErrRecordNotFound
	}

	ctx := context.Background()

	user, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, database.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	v := validator.New()

	if updateParams.Name != nil {
		validateName(v, *updateParams.Name)
		user.Name = *updateParams.Name
	}
	if updateParams.Password != nil {
		validatePasswordPlaintext(v, *updateParams.Password)
	}
	if !v.Valid() {
		return nil, v.GetErrors()
	}

	if updateParams.Password != nil {
		hash, err := hashPassword(*updateParams.Password)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = hash
	}

	result, err := s.queries.UpdateUserById(ctx, store.UpdateUserByIdParams{
		Name:         user.Name,
		PasswordHash: user.PasswordHash,
		ID:           user.ID,
		Version:      user.Version,
	})
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, database.ErrEditConflict
	}

	return &user, nil
}

func (s *UserService) GetForToken(token string) (*store.GetUserFromTokenRow, error) {
	user, err := s.queries.GetUserFromToken(context.Background(), store.GetUserFromTokenParams{
		Hash:   HashToken(token),
		Expiry: time.Now(),
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, database.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
