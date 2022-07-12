package httpsvc

import (
	"net/http"
)

var (
	// ErrUnexpectedInternal is the default unexpected internal server error
	ErrUnexpectedInternal = &HTTPError{
		Status: http.StatusInternalServerError,
		Code:   "internal_error",
		Desc:   "Unexpected internal error occured",
	}
)

// HTTPError represents an error in terms of HTTP info
type HTTPError struct {
	Status int    `json:"_"`
	Code   string `json:"code"`
	Desc   string `json:"description"`
}

// Satisfies go's error interface
func (e HTTPError) Error() string {
	return e.Code + " " + e.Desc
}
