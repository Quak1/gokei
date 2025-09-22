package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Quak1/gokei/internal/auth"
	"github.com/Quak1/gokei/internal/errors"
	"github.com/Quak1/gokei/internal/services"
	"github.com/Quak1/gokei/internal/utils"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: &userService,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	params := services.RegisterUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.ResError(w, err)
		return
	}

	user, err := h.userService.Register(r.Context(), &params)
	if err != nil {
		utils.ResError(w, err)
		return
	}

	utils.ResJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) TokenLogin(w http.ResponseWriter, r *http.Request) {
	params := services.LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.ResError(w, err)
		return
	}

	token, err := h.userService.TokenLogin(r.Context(), &params)
	if err != nil {
		utils.ResError(w, err)
		return
	}

	utils.ResJSON(w, http.StatusOK, token)
}

func (h *UserHandler) EchoUsername(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(*auth.CustomClaims)
	if !ok {
		utils.ResError(w, errors.NewAppError(errors.ErrInternal, "", nil))
		return
	}

	utils.ResJSON(w, http.StatusOK, claims.Username)
}
