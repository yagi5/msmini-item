package item

import (
	"context"
	"errors"
	"testing"
	"time"

	gspanner "cloud.google.com/go/spanner"
	"github.com/google/go-cmp/cmp"
	"github.com/yagi5/msmini-item/domain/data"
	"github.com/yagi5/msmini-item/infrastructure/spanner"
)

func TestSearchByName_Spanner(t *testing.T) {
	var t1 = time.Date(2018, time.January, 1, 12, 30, 0, 0, time.Local)
	var t2 = time.Date(2018, time.January, 2, 12, 30, 0, 0, time.Local)
	var item1 = &data.Item{ID: "id1", Name: "name1", Description: "des1", Price: 100, Category: data.Category("Book"), CreatedAt: t1, UpdatedAt: t2}
	var item2 = &data.Item{ID: "id2", Name: "name2", Description: "des2", Price: 200, Category: data.Category("Book"), CreatedAt: t1, UpdatedAt: t2}
	var item3 = &data.Item{ID: "id3", Name: "name3", Description: "des3", Price: 300, Category: data.Category("Book"), CreatedAt: t1, UpdatedAt: t2}
	var row1, _ = gspanner.NewRow([]string{"id", "name", "description", "price", "category", "createdAt", "updatedAt"}, []interface{}{"id1", "name1", "des1", 100, "Book", t1, t2})
	var row2, _ = gspanner.NewRow([]string{"id", "name", "description", "price", "category", "createdAt", "updatedAt"}, []interface{}{"id2", "name2", "des2", 200, "Book", t1, t2})
	var row3, _ = gspanner.NewRow([]string{"id", "name", "description", "price", "category", "createdAt", "updatedAt"}, []interface{}{"id3", "name3", "des3", 300, "Book", t1, t2})
	var invalidRow1, _ = gspanner.NewRow([]string{"id", "name", "description", "price", "category", "createdAt", "updatedAt"}, []interface{}{nil, "name1", "des1", 100, "Book", t1, t2})
	var invalidRow2, _ = gspanner.NewRow([]string{"id", "name", "description", "price", "category", "createdAt", "updatedAt"}, []interface{}{"id1", nil, "des1", 100, "Book", t1, t2})
	var invalidRow3, _ = gspanner.NewRow([]string{"id", "name", "description", "price", "category", "createdAt", "updatedAt"}, []interface{}{"id1", "name1", nil, 100, "Book", t1, t2})
	var invalidRow4, _ = gspanner.NewRow([]string{"id", "name", "description", "price", "category", "createdAt", "updatedAt"}, []interface{}{"id1", "name1", "des1", nil, "Book", t1, t2})
	var invalidRow5, _ = gspanner.NewRow([]string{"id", "name", "description", "price", "category", "createdAt", "updatedAt"}, []interface{}{"id1", "name1", "des1", 100, nil, t1, t2})
	var invalidRow6, _ = gspanner.NewRow([]string{"id", "name", "description", "price", "category", "createdAt", "updatedAt"}, []interface{}{"id1", "name1", "des1", 100, "Book", nil, t2})
	var invalidRow7, _ = gspanner.NewRow([]string{"id", "name", "description", "price", "category", "createdAt", "updatedAt"}, []interface{}{"id1", "name1", "des1", 100, "Book", t1, nil})
	var tests = []struct {
		name      string
		itemName  string
		limit     int64
		mockQuery func(context.Context, gspanner.Statement) ([]*gspanner.Row, error)
		expected  []*data.Item
		wantErr   bool
	}{
		{
			name:     "validate query and option",
			itemName: "go programming language",
			limit:    20,
			mockQuery: func(ctx context.Context, stmt gspanner.Statement) ([]*gspanner.Row, error) {
				if stmt.SQL != "SELECT * FROM ITEMS WHERE STARTS_WITH(Name, @name) LIMIT @limit" {
					t.Fatalf("SQL invalid: %v", stmt.SQL)
				}
				if diff := cmp.Diff(stmt.Params, map[string]interface{}{"name": "go programming language", "limit": int64(20)}); diff != "" {
					t.Fatalf("params invalid: %v", diff)
				}
				return []*gspanner.Row{row1, row2, row3}, nil
			},
			wantErr:  false,
			expected: []*data.Item{item1, item2, item3},
		},
		{
			name: "return error",
			mockQuery: func(ctx context.Context, stmt gspanner.Statement) ([]*gspanner.Row, error) {
				return nil, errors.New("")
			},
			wantErr: true,
		},
		{
			name: "return error",
			mockQuery: func(ctx context.Context, stmt gspanner.Statement) ([]*gspanner.Row, error) {
				return []*gspanner.Row{invalidRow1}, nil
			},
			wantErr: true,
		},
		{
			name: "return error",
			mockQuery: func(ctx context.Context, stmt gspanner.Statement) ([]*gspanner.Row, error) {
				return []*gspanner.Row{invalidRow2}, nil
			},
			wantErr: true,
		},
		{
			name: "return error",
			mockQuery: func(ctx context.Context, stmt gspanner.Statement) ([]*gspanner.Row, error) {
				return []*gspanner.Row{invalidRow3}, nil
			},
			wantErr: true,
		},
		{
			name: "return error",
			mockQuery: func(ctx context.Context, stmt gspanner.Statement) ([]*gspanner.Row, error) {
				return []*gspanner.Row{invalidRow4}, nil
			},
			wantErr: true,
		},
		{
			name: "return error",
			mockQuery: func(ctx context.Context, stmt gspanner.Statement) ([]*gspanner.Row, error) {
				return []*gspanner.Row{invalidRow5}, nil
			},
			wantErr: true,
		},
		{
			name: "return error",
			mockQuery: func(ctx context.Context, stmt gspanner.Statement) ([]*gspanner.Row, error) {
				return []*gspanner.Row{invalidRow6}, nil
			},
			wantErr: true,
		},
		{
			name: "return error",
			mockQuery: func(ctx context.Context, stmt gspanner.Statement) ([]*gspanner.Row, error) {
				return []*gspanner.Row{invalidRow7}, nil
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ms := &spanner.Mock{MockQuery: tt.mockQuery}
			c := spannerClient{spanner: ms}
			items, err := c.SearchByName(context.Background(), tt.itemName, tt.limit)
			if (err != nil) != tt.wantErr {
				if err != nil {
					t.Fatalf("failed err: %v", err)
				}
			}
			if diff := cmp.Diff(items, tt.expected); diff != "" {
				t.Errorf("failed %v", diff)
			}
		})
	}
}

func TestSearchByName_Dummy(t *testing.T) {
	c := NewDummyClient()
	items, err := c.SearchByName(context.TODO(), "", 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 300 {
		t.Fatal("items is invalid")
	}
}
