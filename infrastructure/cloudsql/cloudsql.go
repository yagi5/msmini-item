package cloudsql

import (
	"context"
	"fmt"

	// blank
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/yagi5/msmini-item/domain/data"
	"github.com/yagi5/msmini-item/domain/repository"
)

// Client is MySQL Client
type Client struct {
	db *sqlx.DB
}

// New returns client
func New(u, pw, host, port, db string) (repository.Item, error) {
	d, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", u, pw, host, port, db))
	if err != nil {
		return nil, errors.Wrap(err, "could not open mysql connection")
	}
	return &Client{db: d}, nil
}

// SearchByName returns searched result
func (c *Client) SearchByName(ctx context.Context, name string, limit int64) ([]*data.Item, error) {
	return nil, nil
}
