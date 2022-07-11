package httpsvc

import (
	"log"
	"net/http"
)

// ErrHandlerFunc is a convenience wrapper around http.HandlerFunc that will handle reporting the error
type ErrHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func handleErrHF(hf ErrHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := hf(w, r)

		if err == nil {
			return
		}

		if herr, ok := err.(*HTTPError); ok {
			if herr.Status < http.StatusInternalServerError || herr.Status == http.StatusServiceUnavailable {
				WriteJSON(r.Context(), w, herr)
				return
			}
		}

		WriteJSON(r.Context(), w, err)
		log.Println(err) // TODO: Deal with this
	}
}
