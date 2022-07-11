package httpsvc

import (
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
)

const (
	toolPrefix = "/_"
)

// Router is a wrapper around chi.Router to not expose chi to main code
type Router struct {
	r chi.Router
}

// Get is used for HTTP GET
func (r *Router) Get(pattern string, hf ErrHandlerFunc) {
	r.r.Get(pattern, handleErrHF(hf))
}

// Post is used for HTTP POST
func (r *Router) Post(pattern string, hf ErrHandlerFunc) {
	r.r.Post(pattern, handleErrHF(hf))
}

// Handle and HandleFunc adds routes for `pattern` that matches all HTTP methods.
func (r *Router) Handle(pattern string, h http.Handler) {
	r.r.Handle(pattern, h)
}

// WithMiddlewares allows injecting middlewares to the stack
func (r *Router) WithMiddlewares(middlewares ...func(http.Handler) http.Handler) {
	r.r.Use(middlewares...)
}

// Group groups a set of routes into one set
func (r *Router) Group(fn func(r *Router)) {
	r.r.Group(func(r chi.Router) {
		fn(&Router{r: r})
	})
}

func prepareHandler(readinessHandler ErrHandlerFunc, routes func(*Router)) http.Handler {
	r := chi.NewRouter()

	r.Get(toolPrefix+"/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		_, _ = fmt.Fprintln(w, "ok")
	})
	r.Get(toolPrefix+"/ready", handleErrHF(readinessHandler))
	pprofRoutes(r)

	r.Group(func(r chi.Router) {
		routes(&Router{r: r})
	})

	return r
}

func pprofRoutes(r chi.Router) {
	prefix := toolPrefix + "/profile"
	r.HandleFunc(prefix+"/*", pprof.Index)
	r.HandleFunc(prefix+"/cmdline", pprof.Cmdline)
	r.HandleFunc(prefix+"/profile", pprof.Profile)
	r.HandleFunc(prefix+"/symbol", pprof.Symbol)
	r.HandleFunc(prefix+"/trace", pprof.Trace)
	r.Handle(prefix+"/goroutine", pprof.Handler("goroutine"))
	r.Handle(prefix+"/threadcreate", pprof.Handler("threadcreate"))
	r.Handle(prefix+"/mutex", pprof.Handler("mutex"))
	r.Handle(prefix+"/heap", pprof.Handler("heap"))
	r.Handle(prefix+"/block", pprof.Handler("block"))
	r.Handle(prefix+"/allocs", pprof.Handler("allocs"))
}
