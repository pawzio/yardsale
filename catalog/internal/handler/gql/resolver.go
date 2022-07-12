//go:generate go run github.com/99designs/gqlgen generate

package gql

import (
	"github.com/99designs/gqlgen/graphql"
)

// NewSchema returns the GraphQL schema
func NewSchema() graphql.ExecutableSchema {
	cfg := Config{
		Resolvers: &resolver{},
	}

	return NewExecutableSchema(cfg)
}

type resolver struct {
}

// Query returns the QueryResolver
func (r *resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct {
	*resolver
}
