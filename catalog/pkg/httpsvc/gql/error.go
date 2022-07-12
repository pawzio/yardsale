package gql

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/errcode"
	"github.com/pawzio/yardsale/catalog/pkg/httpsvc"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

var errLogWrapper = func(ctx context.Context, err error, desc string) {
	log.Println(fmt.Errorf("%s. %w", desc, err)) // TODO: Deal with this
}

// errorPresenter converts Error or error to a presentable format and also handles reporting of errors.
// note: If in prod, gql schema related data is not to be exposed, then this func will need to add redaction.
func errorPresenter() func(ctx context.Context, err error) *gqlerror.Error {
	return func(ctx context.Context, err error) *gqlerror.Error {
		if err == nil {
			return nil
		}

		gqlErr := graphql.DefaultErrorPresenter(ctx, err)

		var herr *httpsvc.HTTPError
		if !errors.As(err, &herr) {
			if gqlErr.Message == "introspection disabled" {
				herr = &httpsvc.HTTPError{
					Status: http.StatusBadRequest,
					Code:   "introspection_disabled",
					Desc:   "Introspection Disabled",
				}
			} else if gqlErr.Extensions != nil {
				herr = parseGQLError(gqlErr)
			} else {
				herr = httpsvc.ErrUnexpectedInternal
			}
		}

		if herr.Status >= http.StatusInternalServerError && herr.Status != http.StatusServiceUnavailable {
			errLogWrapper(ctx, err, "[errorPresenter] encountered unexpected error")
			herr.Desc = httpsvc.ErrUnexpectedInternal.Desc
		}

		gqlErr.Message = herr.Desc
		gqlErr.Extensions = map[string]interface{}{
			"code":        herr.Code,
			"description": herr.Desc,
		}

		return gqlErr
	}
}

func parseGQLError(gqlErr *gqlerror.Error) *httpsvc.HTTPError {
	switch gqlErr.Extensions["code"] {
	case errcode.ValidationFailed:
		herr := &httpsvc.HTTPError{
			Desc:   gqlErr.Message,
			Code:   "gql_validation_failed",
			Status: http.StatusBadRequest,
		}
		return herr
	case errcode.ParseFailed:
		herr := &httpsvc.HTTPError{
			Desc:   gqlErr.Message,
			Code:   "gql_parse_failed",
			Status: http.StatusBadRequest,
		}
		return herr
	default:
		return httpsvc.ErrUnexpectedInternal
	}
}
