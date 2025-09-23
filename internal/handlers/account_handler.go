package handlers

import (
	"net/http"

	"github.com/Quak1/gokei/internal/services"
	"github.com/Quak1/gokei/internal/utils"
)

type AccountHandler struct {
	*BaseHandler
	accountService *services.AccountService
}

func NewAccountHandler(accountService services.AccountService) *AccountHandler {
	return &AccountHandler{
		BaseHandler:    NewBaseHandler(),
		accountService: &accountService,
	}
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	params := services.CreateAccountRequest{}
	if err := h.ParseRequest(w, r, &params); err != nil {
		utils.ResError(w, err)
		return
	}

	claims, err := h.GetJWTClaims(w, r)
	if err != nil {
		utils.ResError(w, err)
		return
	}

	account, err := h.accountService.CreateAccount(r.Context(), claims.Subject, &params)
	if err != nil {
		utils.ResError(w, err)
		return
	}

	utils.ResJSON(w, http.StatusCreated, account)
}
