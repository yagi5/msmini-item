package repository

import (
	"context"

	"github.com/yagi5/msmini-item/domain/data"
)

// Item is accessor to Item data
type Item interface {
	SearchByName(context.Context, string, int64) ([]*data.Item, error)
}
