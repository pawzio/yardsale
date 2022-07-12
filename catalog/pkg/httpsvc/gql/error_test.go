package gql

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/errcode"
	"github.com/pawzio/yardsale/catalog/pkg/httpsvc"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func TestErrorPresenter(t *testing.T) {
	type arg struct {
		givenErr  error
		expErr    *gqlerror.Error
		expLogged bool
	}
	testCases := map[string]arg{
		"go generic err": {
			givenErr: errors.New("some error"),
			expErr: &gqlerror.Error{
				Message: httpsvc.ErrUnexpectedInternal.Desc,
				Extensions: map[string]interface{}{
					"code":        httpsvc.ErrUnexpectedInternal.Code,
					"description": httpsvc.ErrUnexpectedInternal.Desc,
				},
			},
			expLogged: true,
		},
		"http err 500": {
			givenErr: &httpsvc.HTTPError{Status: http.StatusInternalServerError, Code: "some_code", Desc: "some web error"},
			expErr: &gqlerror.Error{
				Message: httpsvc.ErrUnexpectedInternal.Desc,
				Extensions: map[string]interface{}{
					"code":        "some_code",
					"description": httpsvc.ErrUnexpectedInternal.Desc,
				},
			},
			expLogged: true,
		},
		"http err 503": {
			givenErr: &httpsvc.HTTPError{Status: http.StatusServiceUnavailable, Code: "service_unavailable", Desc: "down for maintenance"},
			expErr: &gqlerror.Error{
				Message: "down for maintenance",
				Extensions: map[string]interface{}{
					"code":        "service_unavailable",
					"description": "down for maintenance",
				},
			},
			expLogged: false,
		},
		"http err 400": {
			givenErr: &httpsvc.HTTPError{Status: http.StatusBadRequest, Code: "some_code", Desc: "some web error"},
			expErr: &gqlerror.Error{
				Message: "some web error",
				Extensions: map[string]interface{}{
					"code":        "some_code",
					"description": "some web error",
				},
			},
		},
		"gql errors added via AddError() in vektah/gqlparser/validator/rules": {
			givenErr: &gqlerror.Error{
				Message:    "input: thrown from vektah/gqlparser/validator/rules",
				Extensions: map[string]interface{}{"code": errcode.ValidationFailed},
			},
			expErr: &gqlerror.Error{
				Message: "input: thrown from vektah/gqlparser/validator/rules",
				Extensions: map[string]interface{}{
					"code":        "gql_validation_failed",
					"description": "input: thrown from vektah/gqlparser/validator/rules",
				},
			},
		},
		"gql errors for enums generated in modelgenerated.go": {
			givenErr: &gqlerror.Error{
				Message:    "input: YourInput is not a valid InputType",
				Extensions: map[string]interface{}{"code": errcode.ParseFailed},
			},
			expErr: &gqlerror.Error{
				Message: "input: YourInput is not a valid InputType",
				Extensions: map[string]interface{}{
					"code":        "gql_parse_failed",
					"description": "input: YourInput is not a valid InputType",
				},
			},
		},
		"show location, paths in gqlerror.Error if introspection enabled": {
			givenErr: &gqlerror.Error{
				Message:   "some known error with known path",
				Path:      ast.Path{},
				Locations: []gqlerror.Location{{Line: 1, Column: 5}},
			},
			expErr: &gqlerror.Error{
				Message:   httpsvc.ErrUnexpectedInternal.Desc,
				Path:      ast.Path{},
				Locations: []gqlerror.Location{{Line: 1, Column: 5}},
				Extensions: map[string]interface{}{
					"code":        httpsvc.ErrUnexpectedInternal.Code,
					"description": httpsvc.ErrUnexpectedInternal.Desc,
				},
			},
			expLogged: true,
		},
	}
	for s, tc := range testCases {
		t.Run(s, func(t *testing.T) {
			// Given:
			ctx := graphql.WithFieldContext(context.Background(), &graphql.FieldContext{})

			isLogged := false
			originalErrMonitorWrapper := errLogWrapper
			errLogWrapper = func(ctx context.Context, err error, desc string) {
				isLogged = true
			}
			defer func() { errLogWrapper = originalErrMonitorWrapper }()

			// When:
			result := errorPresenter()(ctx, graphql.ErrorOnPath(ctx, tc.givenErr))

			// Then:
			require.Equal(t, tc.expErr.Message, result.Message)
			require.Equal(t, tc.expErr.Extensions, result.Extensions)
			require.Equal(t, tc.expErr.Locations, result.Locations)
			require.Equal(t, tc.expErr.Path, result.Path)
			require.Equal(t, tc.expLogged, isLogged)
		})
	}
}
