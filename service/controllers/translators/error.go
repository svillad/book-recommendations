package translators

import (
	"encoding/json"
	"net/http"

	"github.com/book-recommendations/service/models"
)

const (
	ErrBadRequest = "invalid query parameters"
)

// ParseError
func ParseError(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	var message string
	switch code {
	case http.StatusBadRequest:
		message = ErrBadRequest
	default:
		message = http.StatusText(code)
	}

	err := models.ResponseError{
		Message: message,
	}
	json.NewEncoder(w).Encode(err)
}
