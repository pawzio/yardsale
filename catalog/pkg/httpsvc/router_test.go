package httpsvc

import (
	"net/http"
	"sort"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestPrepareHandler(t *testing.T) {
	r, err := prepareHandler(func(http.ResponseWriter, *http.Request) error {
		return nil
	}, func(r *Router) {
		r.Get("/get", nil)
		r.Get("/g/{id}", nil)
		r.Post("/post", nil)
		r.Post("/p/{id}", nil)
		r.Group(func(r *Router) {
			r.Get("/abcd/get", nil)
		})
	})
	require.NoError(t, err)

	var routesFound []string
	err = chi.Walk(r, func(method string, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		routesFound = append(routesFound, method+" "+route) // TODO: Add middleware stack check also if possible.
		return nil
	})
	require.NoError(t, err)

	exp := []string{
		"GET /_/healthz",
		"GET /_/ready",

		"CONNECT /_/profile/*", "CONNECT /_/profile/allocs", "CONNECT /_/profile/block", "CONNECT /_/profile/cmdline", "CONNECT /_/profile/goroutine", "CONNECT /_/profile/heap", "CONNECT /_/profile/mutex", "CONNECT /_/profile/profile", "CONNECT /_/profile/symbol", "CONNECT /_/profile/threadcreate", "CONNECT /_/profile/trace",
		"DELETE /_/profile/*", "DELETE /_/profile/allocs", "DELETE /_/profile/block", "DELETE /_/profile/cmdline", "DELETE /_/profile/goroutine", "DELETE /_/profile/heap", "DELETE /_/profile/mutex", "DELETE /_/profile/profile", "DELETE /_/profile/symbol", "DELETE /_/profile/threadcreate", "DELETE /_/profile/trace",
		"GET /_/profile/*", "GET /_/profile/allocs", "GET /_/profile/block", "GET /_/profile/cmdline", "GET /_/profile/goroutine", "GET /_/profile/heap", "GET /_/profile/mutex", "GET /_/profile/profile", "GET /_/profile/symbol", "GET /_/profile/threadcreate", "GET /_/profile/trace",
		"HEAD /_/profile/*", "HEAD /_/profile/allocs", "HEAD /_/profile/block", "HEAD /_/profile/cmdline", "HEAD /_/profile/goroutine", "HEAD /_/profile/heap", "HEAD /_/profile/mutex", "HEAD /_/profile/profile", "HEAD /_/profile/symbol", "HEAD /_/profile/threadcreate", "HEAD /_/profile/trace",
		"OPTIONS /_/profile/*", "OPTIONS /_/profile/allocs", "OPTIONS /_/profile/block", "OPTIONS /_/profile/cmdline", "OPTIONS /_/profile/goroutine", "OPTIONS /_/profile/heap", "OPTIONS /_/profile/mutex", "OPTIONS /_/profile/profile", "OPTIONS /_/profile/symbol", "OPTIONS /_/profile/threadcreate", "OPTIONS /_/profile/trace",
		"PATCH /_/profile/*", "PATCH /_/profile/allocs", "PATCH /_/profile/block", "PATCH /_/profile/cmdline", "PATCH /_/profile/goroutine", "PATCH /_/profile/heap", "PATCH /_/profile/mutex", "PATCH /_/profile/profile", "PATCH /_/profile/symbol", "PATCH /_/profile/threadcreate", "PATCH /_/profile/trace",
		"POST /_/profile/*", "POST /_/profile/allocs", "POST /_/profile/block", "POST /_/profile/cmdline", "POST /_/profile/goroutine", "POST /_/profile/heap", "POST /_/profile/mutex", "POST /_/profile/profile", "POST /_/profile/symbol", "POST /_/profile/threadcreate", "POST /_/profile/trace",
		"PUT /_/profile/*", "PUT /_/profile/allocs", "PUT /_/profile/block", "PUT /_/profile/cmdline", "PUT /_/profile/goroutine", "PUT /_/profile/heap", "PUT /_/profile/mutex", "PUT /_/profile/profile", "PUT /_/profile/symbol", "PUT /_/profile/threadcreate", "PUT /_/profile/trace",
		"TRACE /_/profile/*", "TRACE /_/profile/allocs", "TRACE /_/profile/block", "TRACE /_/profile/cmdline", "TRACE /_/profile/goroutine", "TRACE /_/profile/heap", "TRACE /_/profile/mutex", "TRACE /_/profile/profile", "TRACE /_/profile/symbol", "TRACE /_/profile/threadcreate", "TRACE /_/profile/trace",

		"GET /get",
		"GET /g/{id}",
		"POST /post",
		"POST /p/{id}",
		"GET /abcd/get",
	}
	sort.Strings(exp)

	sort.Strings(routesFound)

	require.EqualValues(t, exp, routesFound)
}
