package repository

import (
	"context"

	"github.com/yagi5/msmini-item/domain/data"
)

// ItemMock is mock
type ItemMock struct {
	SearchByNameMock func(context.Context, string, int64) ([]*data.Item, error)
}

// SearchByName executes injected function for test
func (m *ItemMock) SearchByName(ctx context.Context, name string, limit int64) ([]*data.Item, error) {
	return m.SearchByNameMock(ctx, name, limit)
}
