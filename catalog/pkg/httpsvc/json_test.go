package httpsvc

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteJSON(t *testing.T) {
	type testCase struct {
		givenObj  interface{}
		expStatus int
	}
	tcs := map[string]testCase{
		"no payload": {expStatus: http.StatusOK},
		"payload": {
			givenObj: map[string]interface{}{
				"k1": "v1",
				"k2": 2.3,
				"k3": 123.554,
				"k4": map[string]interface{}{
					"c1": "v1",
					"c2": []interface{}{"1", "2"},
				},
			},
			expStatus: http.StatusOK,
		},
		"400 error": {
			givenObj:  &HTTPError{Status: http.StatusBadRequest, Code: "err", Desc: "desc"},
			expStatus: http.StatusBadRequest,
		},
		"401 error": {
			givenObj:  &HTTPError{Status: http.StatusUnauthorized, Code: "err", Desc: "desc"},
			expStatus: http.StatusUnauthorized,
		},
		"500 err": {
			givenObj:  &HTTPError{Status: http.StatusBadRequest, Code: "err", Desc: "desc"},
			expStatus: http.StatusBadRequest,
		},
	}
	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			// Given:
			w := httptest.NewRecorder()

			// When:
			WriteJSON(context.Background(), w, tc.givenObj)

			// Then:
			require.Equal(t, tc.expStatus, w.Result().StatusCode)
			require.Equal(t, "application/json", w.Result().Header.Get("Content-Type"))
			bytes, err := io.ReadAll(w.Result().Body)
			defer w.Result().Body.Close()
			require.NoError(t, err)

			if tc.givenObj == nil {
				resp := map[string]interface{}{}
				err = json.Unmarshal(bytes, &resp)
				require.NoError(t, err)
				require.Equal(t, map[string]interface{}(nil), resp)
				return
			}

			if v, ok := tc.givenObj.(*HTTPError); ok {
				var resp HTTPError
				err = json.Unmarshal(bytes, &resp)
				require.NoError(t, err)

				if v.Status == http.StatusInternalServerError {
					require.Equal(t, ErrUnexpectedInternal.Code, resp.Code)
					require.Equal(t, ErrUnexpectedInternal.Desc, resp.Desc)
					return
				}

				require.Equal(t, v.Code, resp.Code)
				require.Equal(t, v.Desc, resp.Desc)
				return
			}

			resp := map[string]interface{}{}
			err = json.Unmarshal(bytes, &resp)
			require.NoError(t, err)
			require.EqualValues(t, tc.givenObj, resp)
		})
	}
}
