package item

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/yagi5/msmini-item/domain/data"
	"github.com/yagi5/msmini-item/infrastructure/spanner"
)

// Client contains spanner client
type Client struct {
	spanner spanner.Spanner
}

// SearchByName returns searched result
func (c *Client) SearchByName(ctx context.Context, name string, limit int64) ([]*data.Item, error) {
	sql := "SELECT * FROM ITEMS WHERE STARTS_WITH(Name, @name) LIMIT @limit"
	params := map[string]interface{}{"name": name, "limit": limit}
	stmt := spanner.NewStatement(sql, params)
	rows, err := c.spanner.Query(ctx, stmt)
	if err != nil {
		return nil, errors.Wrap(err, "Query failed")
	}
	var items []*data.Item
	for _, row := range rows {
		var id string
		if err := row.Column(0, &id); err != nil {
			return nil, errors.Wrap(err, "scan id failed")
		}
		var name string
		if err := row.Column(1, &name); err != nil {
			return nil, errors.Wrap(err, "scan name failed")
		}
		var description string
		if err := row.Column(2, &description); err != nil {
			return nil, errors.Wrap(err, "scan description failed")
		}
		var price int64
		if err := row.Column(3, &price); err != nil {
			return nil, errors.Wrap(err, "scan price failed")
		}
		var category string
		if err := row.Column(4, &category); err != nil {
			return nil, errors.Wrap(err, "scan category failed")
		}
		var createdAt time.Time
		if err := row.Column(5, &createdAt); err != nil {
			return nil, errors.Wrap(err, "scan createdAt failed")
		}
		var updatedAt time.Time
		if err := row.Column(6, &updatedAt); err != nil {
			return nil, errors.Wrap(err, "scan updatedAt failed")
		}

		item := &data.Item{
			ID:          id,
			Name:        name,
			Description: description,
			Price:       price,
			Category:    data.Category(category),
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}
		items = append(items, item)
	}
	return items, nil
}
