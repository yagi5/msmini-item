package usecase

import (
	"context"

	"github.com/yagi5/msmini-item/domain/data"
	"github.com/yagi5/msmini-item/domain/repository"
	"github.com/yagi5/msmini-item/log"
)

// Item is item usecase
type Item interface {
	Search(context.Context, *ItemSearchInput) ([]*data.Item, error)
}

// Item usecase handler
type item struct {
	itemRepo repository.Item
	logger   *log.Logger
}

// New returns item usecase
func New(logger *log.Logger, r repository.Item) Item {
	return item{logger: logger, itemRepo: r}
}

// Search returns items queried by input
func (u item) Search(ctx context.Context, in *ItemSearchInput) ([]*data.Item, error) {
	u.logger.Info("search start", log.F("name", in.Name))
	if in.Limit == 0 {
		in.Limit = 20 // Default
	}
	return u.itemRepo.SearchByName(ctx, in.Name, in.Limit)
}
