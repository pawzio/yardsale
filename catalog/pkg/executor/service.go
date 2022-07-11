package executor

import (
	"context"
)

// Service is an executable service that the executor can run
type Service func(ctx context.Context) error
