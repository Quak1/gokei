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

type AccountHandler struct {
	accountService *service.AccountService
}

func NewAccountHandler(svc *service.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: svc,
	}
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Type           store.AccountType `json:"type"`
		Name           string            `json:"name"`
		InitialBalance int64             `json:"initial_balance"`
	}

	err := response.ReadJSON(w, r, &input)
	if err != nil {
		response.BadRequestResponse(w, r, err)
		return
	}

	ctxUser := appcontext.GetContextUser(r)

	params := store.CreateAccountParams{
		Type:         input.Type,
		Name:         input.Name,
		UserID:       ctxUser.ID,
		BalanceCents: input.InitialBalance,
	}

	account, err := h.accountService.Create(&params)
	if err != nil {
		var validationErr *validator.ValidationError

		switch {
		case errors.As(err, &validationErr):
			response.FailedValidationResponse(w, r, validationErr)
		default:
			response.ServerErrorResponse(w, r, err)
		}

		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/accounts/%d", account.ID))

	err = response.Created(w, response.Envelope{"account": account}, headers)
	if err != nil {
		response.ServerErrorResponse(w, r, err)
	}
}

func (h *AccountHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctxUser := appcontext.GetContextUser(r)
	accounts, err := h.accountService.GetAll(ctxUser.ID)

	if err != nil {
		response.ServerErrorResponse(w, r, err)
	}

	err = response.OK(w, response.Envelope{"accounts": accounts})
	if err != nil {
		response.ServerErrorResponse(w, r, err)
	}
}

func (h *AccountHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := readIntParam(r, "accountID")
	if err != nil {
		response.BadRequestResponseGeneric(w, r)
		return
	}

	ctxUser := appcontext.GetContextUser(r)

	account, err := h.accountService.GetByID(int32(id), ctxUser.ID)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrRecordNotFound):
			response.NotFoundResponse(w, r)
		default:
			response.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = response.OK(w, response.Envelope{"account": account})
	if err != nil {
		response.ServerErrorResponse(w, r, err)
	}
}

func (h *AccountHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id, err := readIntParam(r, "accountID")
	if err != nil {
		response.BadRequestResponseGeneric(w, r)
		return
	}

	ctxUser := appcontext.GetContextUser(r)

	err = h.accountService.DeleteByID(int32(id), ctxUser.ID)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrRecordNotFound):
			response.NotFoundResponse(w, r)
		default:
			response.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = response.OK(w, response.Envelope{"message": "account successfully deleted"})
	if err != nil {
		response.ServerErrorResponse(w, r, err)
	}
}

func (h *AccountHandler) GetSumBalance(w http.ResponseWriter, r *http.Request) {
	id, err := readIntParam(r, "accountID")
	if err != nil {
		response.BadRequestResponseGeneric(w, r)
		return
	}

	ctxUser := appcontext.GetContextUser(r)

	balance, err := h.accountService.GetSumBalance(int32(id), ctxUser.ID)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrRecordNotFound):
			response.NotFoundResponse(w, r)
		default:
			response.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = response.OK(w, response.Envelope{"balance_cents": balance})
	if err != nil {
		response.ServerErrorResponse(w, r, err)
	}
}

func (h *AccountHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	id, err := readIntParam(r, "accountID")
	if err != nil {
		response.BadRequestResponseGeneric(w, r)
		return
	}

	var input service.UpdateAccountParams
	err = response.ReadJSON(w, r, &input)
	if err != nil {
		response.BadRequestResponse(w, r, err)
		return
	}

	ctxUser := appcontext.GetContextUser(r)

	account, err := h.accountService.UpdateByID(int32(id), ctxUser.ID, &input)
	if err != nil {
		var validationErr *validator.ValidationError
		switch {
		case errors.As(err, &validationErr):
			response.FailedValidationResponse(w, r, validationErr)
		case errors.Is(err, database.ErrRecordNotFound):
			response.NotFoundResponse(w, r)
		case errors.Is(err, database.ErrEditConflict):
			response.ConflictResponse(w, r)
		default:
			response.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = response.OK(w, response.Envelope{"account": account})
	if err != nil {
		response.ServerErrorResponse(w, r, err)
	}
}
