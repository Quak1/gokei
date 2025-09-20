package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Quak1/gokei/internal/errors"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type Respose struct {
	Data any `json:"data"`
}

func ResError(w http.ResponseWriter, err error) {
	var statusCode int
	var errRes ErrorResponse

	if appErr, ok := err.(*errors.AppError); ok {
		statusCode = appErr.StatusCode
		errRes = ErrorResponse{Error: appErr.UserMessage}
	} else {
		statusCode = http.StatusInternalServerError
		errRes = ErrorResponse{Error: "An internal error ocurred"}
	}

	log.Println(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if encodeErr := json.NewEncoder(w).Encode(errRes); encodeErr != nil {
		log.Printf("Error encoding JSON response: %s", encodeErr)
	}
}

func ResJSON(w http.ResponseWriter, code int, payload any) {
	res := Respose{Data: payload}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if encodeErr := json.NewEncoder(w).Encode(res); encodeErr != nil {
		log.Printf("Error encoding JSON response: %s", encodeErr)
	}
}
