package httpsvc

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrHandlerFunc(t *testing.T) {
	type arg struct {
		givenErr  error
		expStatus int
		expErr    *HTTPError
	}
	tcs := map[string]arg{
		"no err": {
			expStatus: http.StatusOK,
		},
		"go generic err": {
			givenErr:  errors.New("some error"),
			expStatus: http.StatusInternalServerError,
			expErr:    ErrUnexpectedInternal,
		},
		"http err 500": {
			givenErr:  &HTTPError{Status: http.StatusInternalServerError, Code: "code", Desc: "desc"},
			expStatus: http.StatusInternalServerError,
			expErr:    &HTTPError{Status: http.StatusInternalServerError, Code: ErrUnexpectedInternal.Code, Desc: ErrUnexpectedInternal.Desc},
		},
		"http err 503": {
			givenErr:  &HTTPError{Status: http.StatusServiceUnavailable, Code: "code", Desc: "desc"},
			expStatus: http.StatusServiceUnavailable,
			expErr:    &HTTPError{Code: "code", Desc: "desc"},
		},
		"http err 400": {
			givenErr:  &HTTPError{Status: http.StatusBadRequest, Code: "code", Desc: "desc"},
			expStatus: http.StatusBadRequest,
			expErr:    &HTTPError{Code: "code", Desc: "desc"},
		},
		"http err 401": {
			givenErr:  &HTTPError{Status: http.StatusUnauthorized, Code: "code", Desc: "desc"},
			expStatus: http.StatusUnauthorized,
			expErr:    &HTTPError{Code: "code", Desc: "desc"},
		},
	}
	for s, tc := range tcs {
		t.Run(s, func(t *testing.T) {
			// Given:
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			// When:
			handleErrHF(func(w http.ResponseWriter, r *http.Request) error {
				return tc.givenErr
			}).ServeHTTP(w, r)

			// Then:
			require.Equal(t, tc.expStatus, w.Code)
			bytes, err := io.ReadAll(w.Result().Body)
			defer w.Result().Body.Close()

			var actErr HTTPError
			err = json.Unmarshal(bytes, &actErr)
			if tc.expErr == nil {
				require.Error(t, err)
			} else {
				require.Nil(t, err)
				require.Equal(t, tc.expErr.Code, actErr.Code)
				require.Equal(t, tc.expErr.Desc, actErr.Desc)
			}
		})
	}
}
