package gql

import (
	"context"
)

// NoOp is a No-Op resolver to set up the GQL configs. Once proper query resolvers are added, this will be removed.
func (*queryResolver) NoOp(context.Context) (*bool, error) {
	return nil, nil
}
