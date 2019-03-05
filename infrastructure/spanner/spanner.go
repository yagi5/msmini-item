package spanner

import (
	"context"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

// Spanner is interface for Google Cloud Spanner
type Spanner interface {
	Queryer
}

// Queryer provides Read()
type Queryer interface {
	Query(ctx context.Context, stmt spanner.Statement) ([]*spanner.Row, error)
}

// Client is spanner client
type Client struct {
	client *spanner.Client
}

// New returns spanner client
func New(c *spanner.Client) Spanner {
	return &Client{c}
}

// Query returns query results
// By this, Client satisfies Queryer Interface
func (c *Client) Query(ctx context.Context, stmt spanner.Statement) (rows []*spanner.Row, err error) {
	iter := c.client.ReadOnlyTransaction().Query(ctx, stmt)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}
	return
}

// NewStatement returns spanner.Statement
func NewStatement(sql string, params map[string]interface{}) spanner.Statement {
	return spanner.Statement{SQL: sql, Params: params}
}
