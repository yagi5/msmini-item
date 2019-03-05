package cloudsql

import (
	"context"
	"fmt"

	// blank
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// CloudSQL is GCP CloudSQL interface
type CloudSQL interface {
	Queryer
}

// Queryer provides Read()
type Queryer interface {
	Query(ctx context.Context, query string, args ...interface{}) *sqlx.Row
}

// Client is MySQL Client
type Client struct {
	db *sqlx.DB
}

// New returns client
func New(u, pw, host, port, db string) (CloudSQL, error) {
	d, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", u, pw, host, port, db))
	if err != nil {
		return nil, errors.Wrap(err, "could not open mysql connection")
	}
	return &Client{db: d}, nil
}

// Query implements Queryer
func (c *Client) Query(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return c.db.QueryRowxContext(ctx, query, args)
}
