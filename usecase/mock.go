package usecase

import (
	"context"

	"github.com/yagi5/msmini-item/domain/data"
)

// ItemMock is item usecase mock
type ItemMock struct {
	SearchMock func(context.Context, *ItemSearchInput) ([]*data.Item, error)
}

// Search executes injected mock
func (u *ItemMock) Search(ctx context.Context, in *ItemSearchInput) ([]*data.Item, error) {
	return u.SearchMock(ctx, in)
}
