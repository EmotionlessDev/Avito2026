package common

import (
	"errors"
	"net/http"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/helpers"
)

var (
	ErrInvalidRole   error = errors.New("invalid role")
	ErrNilTx         error = errors.New("nil transaction")
	ErrDuplicateRoom error = errors.New("room with the same name already exists")
	ErrBeginTx       error = errors.New("failed to begin transaction")
)

func errorResponse(w http.ResponseWriter, status int, err error, message interface{}) {
	helpers.WriteJSON(w, status, helpers.Envelope{
		"error":   err.Error(),
		"message": message,
	}, nil)
}

func MethodNotAllowedResponse(w http.ResponseWriter) {
	errorResponse(w, http.StatusMethodNotAllowed, errors.New("method not allowed"), "the requested method is not allowed for the specified resource")
}

func InternalServerErrorResponse(w http.ResponseWriter, err error) {
	errorResponse(w, http.StatusInternalServerError, errors.New("internal server error"), err.Error())
}

func BadRequestResponse(w http.ResponseWriter, err error) {
	errorResponse(w, http.StatusBadRequest, errors.New("bad request"), err.Error())
}

func NotFoundResponse(w http.ResponseWriter, err error) {
	errorResponse(w, http.StatusNotFound, errors.New("not found"), err.Error())
}

func FailedValidationResponse(w http.ResponseWriter, errors map[string]string) {
	helpers.WriteJSON(w, http.StatusUnprocessableEntity, helpers.Envelope{
		"error":   "failed validation",
		"message": errors,
	}, nil)
}
