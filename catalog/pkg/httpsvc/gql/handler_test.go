package gql

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/ast"
)

func TestHandler(t *testing.T) {
	// Given && When:
	h := Handler(&mockExecutableSchema{}, false)

	// Then:
	require.NotNil(t, h)

	// Given && When:
	h = Handler(&mockExecutableSchema{}, true)

	// Then:
	require.NotNil(t, h)
}

type mockExecutableSchema struct{}

func (mockExecutableSchema) Schema() *ast.Schema {
	panic("implement me")
}

func (mockExecutableSchema) Complexity(string, string, int, map[string]interface{}) (int, bool) {
	panic("implement me")
}

func (mockExecutableSchema) Exec(context.Context) graphql.ResponseHandler {
	panic("implement me")
}
