package usecase

import (
	"context"

	"github.com/pkg/errors"
	"github.com/yagi5/msmini-item/domain/data"
	"github.com/yagi5/msmini-item/usecase/repository"
)

// Repositories is structured repositories
type Repositories struct {
	Item repository.Item
}

// ItemUsecase usecase handler
type ItemUsecase struct {
	repos Repositories
}

// Search returns items queried by input
func (u ItemUsecase) Search(ctx context.Context, in *ItemSearchInput) ([]*data.Item, error) {
	if in.Limit == 0 {
		in.Limit = 20 // Default
	}
	if in.Name == "" {
		return nil, errors.New("name is required")
	}
	return u.repos.Item.SearchByName(ctx, in.Name, in.Limit)
}
