//go:build tools
// +build tools

package tools

import (
	_ "github.com/99designs/gqlgen"                       // Need to use go generate for gqlgen
	_ "github.com/99designs/gqlgen/graphql"               // Need to use go generate for gqlgen
	_ "github.com/99designs/gqlgen/graphql/introspection" // Need to use go generate for gqlgen
)
