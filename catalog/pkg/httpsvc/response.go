package httpsvc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// WriteJSON writes JSON output to the response writer
func WriteJSON(_ context.Context, w http.ResponseWriter, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")

	var status int

	switch t := obj.(type) {
	case *HTTPError:
		if t.Status < http.StatusInternalServerError || t.Status == http.StatusServiceUnavailable {
			obj = t
			status = t.Status
		} else {
			obj = ErrUnexpectedInternal
			status = ErrUnexpectedInternal.Status
		}
	case error:
		obj = ErrUnexpectedInternal
	default:
		status = http.StatusOK
	}

	b, err := json.Marshal(obj)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Errorf("writeJSON: json marshal failed. err: %w", err)) // TODO: Deal with this
		return
	}

	w.WriteHeader(status)

	_, _ = w.Write(b)
}
