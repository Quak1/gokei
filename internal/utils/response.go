package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Quak1/gokei/internal/apperrors"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type ValidationErrorResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

type Respose struct {
	Data any `json:"data"`
}

func sendResponse(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if encodeErr := json.NewEncoder(w).Encode(payload); encodeErr != nil {
		log.Printf("Error encoding JSON response: %s", encodeErr)
	}
}

func ResJSON(w http.ResponseWriter, statusCode int, payload any) {
	res := Respose{Data: payload}
	sendResponse(w, statusCode, res)
}

func ResError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		sendResponse(w, appErr.StatusCode(), ErrorResponse{Error: appErr.Message})
	} else if appErr, ok := err.(*apperrors.ValidationError); ok {
		sendResponse(w, appErr.StatusCode(), ValidationErrorResponse{
			Message: appErr.Message,
			Errors:  appErr.Fields,
		})
		return
	} else {
		sendResponse(w, http.StatusInternalServerError, ErrorResponse{Error: "An internal error ocurred"})
	}

	log.Println(err)
}
