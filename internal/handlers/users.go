package handlers

import (
	"encoding/json"
	"net/http"

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
