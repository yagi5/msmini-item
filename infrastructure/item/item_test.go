package item

import (
	"context"
	"testing"
	"time"

	gspanner "cloud.google.com/go/spanner"
	"github.com/google/go-cmp/cmp"
	"github.com/yagi5/msmini-item/domain/data"
	"github.com/yagi5/msmini-item/infrastructure/spanner"
)

func getTestRows() []*gspanner.Row {
	l := "2006-01-02 15:04:05"
	t1, _ := time.Parse(l, "2018-01-01 12:30:00")
	t2, _ := time.Parse(l, "2018-01-02 12:30:00")
	item1, _ := gspanner.NewRow([]string{
		"id", "name", "description", "price", "category", "createdAt", "updatedAt",
	}, []interface{}{
		"id1", "name1", "des1", 100, "Book", t1, t2,
	})
	item2, _ := gspanner.NewRow([]string{
		"id", "name", "description", "price", "category", "createdAt", "updatedAt",
	}, []interface{}{
		"id2", "name2", "des2", 200, "Book", t1, t2,
	})
	item3, _ := gspanner.NewRow([]string{
		"id", "name", "description", "price", "category", "createdAt", "updatedAt",
	}, []interface{}{
		"id3", "name3", "des3", 300, "Book", t1, t2,
	})
	return []*gspanner.Row{item1, item2, item3}
}

func getTestItems() []*data.Item {
	l := "2006-01-02 15:04:05"
	t1, _ := time.Parse(l, "2018-01-01 12:30:00")
	t2, _ := time.Parse(l, "2018-01-02 12:30:00")
	return []*data.Item{
		&data.Item{
			ID:          "id1",
			Name:        "name1",
			Description: "des1",
			Price:       100,
			Category:    data.Category("Book"),
			CreatedAt:   t1,
			UpdatedAt:   t2,
		},
		&data.Item{
			ID:          "id2",
			Name:        "name2",
			Description: "des2",
			Price:       200,
			Category:    data.Category("Book"),
			CreatedAt:   t1,
			UpdatedAt:   t2,
		},
		&data.Item{
			ID:          "id3",
			Name:        "name3",
			Description: "des3",
			Price:       300,
			Category:    data.Category("Book"),
			CreatedAt:   t1,
			UpdatedAt:   t2,
		},
	}
}

func TestSearchByName(t *testing.T) {
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
				return getTestRows(), nil
			},
			wantErr:  false,
			expected: getTestItems(),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ms := &spanner.Mock{MockQuery: tt.mockQuery}
			c := Client{spanner: ms}
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
