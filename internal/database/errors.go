package database

import (
	"errors"
	"fmt"

	"github.com/lib/pq"
)

var (
	ErrRecordNotFound        = errors.New("record not found")
	ErrEditConflict          = errors.New("edit conflict")
	ErrInvalidCategory       = errors.New("This category does not exist")
	ErrUpdateInitialCategory = errors.New("Can't update initial category")
	ErrInvalidAccount        = errors.New("This account does not exist")
)

func HandleForeignKeyError(err error) error {
	if pqErr, ok := err.(*pq.Error); ok {
		if pqErr.Code == "23503" {
			switch pqErr.Constraint {
			case "transactions_category_id_fkey":
				return ErrInvalidCategory
			case "transactions_account_id_fkey":
				return ErrInvalidAccount
			default:
				return fmt.Errorf("reference does not exist: %s", pqErr.Constraint)
			}
		}
	}

	return err
}

func IsUniqueContraintViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505"
	}

	return false
}
