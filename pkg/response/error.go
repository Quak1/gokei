package response

import (
	"fmt"
	"net/http"

	"github.com/Quak1/gokei/pkg/validator"
)

func logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	logger.Error(err.Error(), "method", method, "uri", uri)
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := Envelope{"error": message}

	err := WriteJSON(w, status, env, nil)
	if err != nil {
		logError(r, err)
		w.WriteHeader(500)
	}
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	logError(r, err)

	message := "Something went wrong on our end. Please try againg later."
	ErrorResponse(w, r, http.StatusInternalServerError, message)
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "The resource you're looking for doesn't exist."
	ErrorResponse(w, r, http.StatusNotFound, message)
}

func MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %s method is not supported for this endpoint.", r.Method)
	ErrorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

func BadRequestResponseGeneric(w http.ResponseWriter, r *http.Request) {
	message := "We couldn't understand your request. Please check your input and try again."
	ErrorResponse(w, r, http.StatusBadRequest, message)
}

func FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors *validator.ValidationError) {
	ErrorResponse(w, r, http.StatusUnprocessableEntity, errors.Errors)
}

func UnauthorizedResponse(w http.ResponseWriter, r *http.Request) {
	message := "You need to be logged in to access this resource."
	ErrorResponse(w, r, http.StatusUnauthorized, message)
}

func ForbiddenResponse(w http.ResponseWriter, r *http.Request) {
	message := "You don't have permission to access this resource."
	ErrorResponse(w, r, http.StatusForbidden, message)
}

func ConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "This resource already exists or conflicts with an existing resource."
	ErrorResponse(w, r, http.StatusConflict, message)
}

func RateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "You've made too many requests. Please slow down and try again later."
	ErrorResponse(w, r, http.StatusTooManyRequests, message)
}

func InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "Invalid credentials. Verify your login information."
	ErrorResponse(w, r, http.StatusUnauthorized, message)
}
