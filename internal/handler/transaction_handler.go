package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Quak1/gokei/internal/appcontext"
	"github.com/Quak1/gokei/internal/database"
	"github.com/Quak1/gokei/internal/database/store"
	"github.com/Quak1/gokei/internal/service"
	"github.com/Quak1/gokei/pkg/response"
	"github.com/Quak1/gokei/pkg/validator"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
}

func NewTransactionHandler(svc *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: svc,
	}
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input store.CreateTransactionParams

	err := response.ReadJSON(w, r, &input)
	if err != nil {
		response.BadRequestResponse(w, r, err)
		return
	}

	ctxUser := appcontext.GetContextUser(r)

	transaction, err := h.transactionService.Create(ctxUser.ID, &input)
	if err != nil {
		var validationErr *validator.ValidationError

		switch {
		case errors.As(err, &validationErr):
			response.FailedValidationResponse(w, r, validationErr)
		case errors.Is(err, database.ErrInvalidCategory):
			response.BadRequestResponse(w, r, err)
		case errors.Is(err, database.ErrInvalidAccount), errors.Is(err, database.ErrRecordNotFound):
			response.NotFoundResponse(w, r)
		default:
			response.ServerErrorResponse(w, r, err)
		}

		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/transactions/%d", transaction.ID))

	err = response.Created(w, response.Envelope{"transaction": transaction}, headers)
	if err != nil {
		response.ServerErrorResponse(w, r, err)
	}
}

func (h *TransactionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctxUser := appcontext.GetContextUser(r)

	transactions, err := h.transactionService.GetAll(ctxUser.ID)
	if err != nil {
		response.ServerErrorResponse(w, r, err)
		return
	}

	err = response.OK(w, response.Envelope{"transactions": transactions})
	if err != nil {
		response.ServerErrorResponse(w, r, err)
		return
	}
}

func (h *TransactionHandler) GetAccountTransactions(w http.ResponseWriter, r *http.Request) {
	accountID, err := readIntParam(r, "accountID")
	if err != nil {
		response.BadRequestResponseGeneric(w, r)
		return
	}

	ctxUser := appcontext.GetContextUser(r)

	transactions, err := h.transactionService.GetAllTRansactionsForAccountID(int32(accountID), ctxUser.ID)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrRecordNotFound):
			response.NotFoundResponse(w, r)
		default:
			response.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = response.OK(w, response.Envelope{"transactions": transactions})
	if err != nil {
		response.ServerErrorResponse(w, r, err)
		return
	}
}

func (h *TransactionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := readIntParam(r, "transactionID")
	if err != nil {
		response.BadRequestResponseGeneric(w, r)
		return
	}

	ctxUser := appcontext.GetContextUser(r)

	transaction, err := h.transactionService.GetByID(int32(id), ctxUser.ID)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrRecordNotFound):
			response.NotFoundResponse(w, r)
		default:
			response.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = response.OK(w, response.Envelope{"transaction": transaction})
	if err != nil {
		response.ServerErrorResponse(w, r, err)
	}
}

func (h *TransactionHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id, err := readIntParam(r, "transactionID")
	if err != nil {
		response.BadRequestResponseGeneric(w, r)
		return
	}

	ctxUser := appcontext.GetContextUser(r)

	err = h.transactionService.DeleteByID(int32(id), ctxUser.ID)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrRecordNotFound):
			response.NotFoundResponse(w, r)
		default:
			response.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = response.OK(w, response.Envelope{"message": "transaction successfully deleted"})
	if err != nil {
		response.ServerErrorResponse(w, r, err)
	}
}

func (h *TransactionHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	id, err := readIntParam(r, "transactionID")
	if err != nil {
		response.BadRequestResponseGeneric(w, r)
		return
	}

	var input service.UpdateTransactionParams
	err = response.ReadJSON(w, r, &input)
	if err != nil {
		response.BadRequestResponse(w, r, err)
		return
	}

	ctxUser := appcontext.GetContextUser(r)

	transaction, err := h.transactionService.UpdateByID(int32(id), ctxUser.ID, &input)
	if err != nil {
		var validationErr *validator.ValidationError
		switch {
		case errors.As(err, &validationErr):
			response.FailedValidationResponse(w, r, validationErr)
		case errors.Is(err, database.ErrRecordNotFound):
			response.NotFoundResponse(w, r)
		case errors.Is(err, database.ErrEditConflict):
			response.ConflictResponse(w, r)
		case errors.Is(err, database.ErrInvalidAccount), errors.Is(err, database.ErrInvalidCategory):
			response.BadRequestResponse(w, r, err)
		default:
			response.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = response.OK(w, response.Envelope{"transaction": transaction})
	if err != nil {
		response.ServerErrorResponse(w, r, err)
	}
}
