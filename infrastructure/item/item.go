package item

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/yagi5/msmini-item/domain/data"
	"github.com/yagi5/msmini-item/domain/repository"
	"github.com/yagi5/msmini-item/infrastructure/cloudsql"
	"github.com/yagi5/msmini-item/infrastructure/spanner"
)

// spannerClient contains spanner client
type spannerClient struct {
	spanner spanner.Spanner
}

// csqlClient contains cloudSQL client
type csqlClient struct {
	csql cloudsql.CloudSQL
}

// dummyClient returns dummy data client
type dummyClient struct{}

// NewSpannerClient returns item repository client
func NewSpannerClient(s spanner.Spanner) repository.Item {
	return &spannerClient{s}
}

// NewCloudSQLClient returns item repository client
func NewCloudSQLClient(c cloudsql.CloudSQL) repository.Item {
	return &csqlClient{c}
}

// NewDummyClient returns item repository client
func NewDummyClient() repository.Item {
	return &dummyClient{}
}

// SearchByName returns searched result
func (c *spannerClient) SearchByName(ctx context.Context, name string, limit int64) (items []*data.Item, err error) {
	defer func() {
		if rerr := recover(); rerr != nil {
			log.Printf("recovered: %v", rerr)
			items, err = NewDummyClient().SearchByName(ctx, name, limit)
		}
	}()

	sql := "SELECT * FROM ITEMS LIMIT @limit"
	if name != "" {
		sql = "SELECT * FROM ITEMS WHERE STARTS_WITH(Name, @name) LIMIT @limit"
	}
	params := map[string]interface{}{"name": name, "limit": limit}
	stmt := spanner.NewStatement(sql, params)
	rows, err := c.spanner.Query(ctx, stmt)
	if err != nil {
		return nil, errors.Wrap(err, "Query failed")
	}
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

// SearchByName returns searched result
func (c *csqlClient) SearchByName(ctx context.Context, name string, limit int64) (items []*data.Item, err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("recovered: %v", err)
			items, err = NewDummyClient().SearchByName(ctx, name, limit)
		}
	}()

	query := "SELECT * FROM Items LIMIT $1"
	args := []interface{}{limit}
	if name != "" {
		query = "SELECT * FROM Items WHERE Name LIKE '$1%' LIMIT $2"
		args = []interface{}{name, limit}
	}
	row := c.csql.Query(ctx, query, args)
	err = row.StructScan(items)
	if err != nil {
		return nil, err
	}
	return items, nil
}

// SearchByName returns searched result
func (c *dummyClient) SearchByName(ctx context.Context, name string, limit int64) ([]*data.Item, error) {
	items := []*data.Item{}
	for i := 1; i <= 300; i++ {
		item := &data.Item{
			ID:          strconv.Itoa(i),
			Name:        fmt.Sprintf("dummy item #%d", i),
			Description: fmt.Sprintf("description for dummy item #%d", i),
			Price:       int64(i),
			Category:    data.Book,
			CreatedAt:   time.Date(2019, time.February, 10, 12, 0, 0, i, time.Local),
			UpdatedAt:   time.Date(2019, time.February, 10, 12, 0, 0, i, time.Local),
		}
		items = append(items, item)
	}
	return items, nil
}
