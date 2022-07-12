package main

import (
	"net/http"

	"github.com/pawzio/yardsale/catalog/internal/handler/gql"
	"github.com/pawzio/yardsale/catalog/pkg/httpsvc"
	httpsvcgql "github.com/pawzio/yardsale/catalog/pkg/httpsvc/gql"
)

func routes(r *httpsvc.Router) {
	r.Handle("/graphql", httpsvcgql.Handler(gql.NewSchema(), false)) // TODO: Deal with the introspection
}

func readinessHandler() httpsvc.ErrHandlerFunc { // TODO: Deal with this
	return func(http.ResponseWriter, *http.Request) error {
		return nil
	}
}
