package httperr

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func Write(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func BadRequest(w http.ResponseWriter, message string) {
	Write(w, http.StatusBadRequest, message)
}

func Unauthorized(w http.ResponseWriter, message string) {
	Write(w, http.StatusUnauthorized, message)
}

func NotFound(w http.ResponseWriter, message string) {
	Write(w, http.StatusNotFound, message)
}

func Conflict(w http.ResponseWriter, message string) {
	Write(w, http.StatusConflict, message)
}

func Internal(w http.ResponseWriter, message string) {
	Write(w, http.StatusInternalServerError, message)
}