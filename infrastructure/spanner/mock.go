package spanner

import (
	"context"

	"cloud.google.com/go/spanner"
)

// Mock is mock
type Mock struct {
	MockQuery func(ctx context.Context, stmt spanner.Statement) ([]*spanner.Row, error)
}

// Query returns injected function
func (m *Mock) Query(ctx context.Context, stmt spanner.Statement) ([]*spanner.Row, error) {
	return m.MockQuery(ctx, stmt)
}
