package data_test

import (
	"testing"
	"time"

	"github.com/yagi5/msmini-item/domain/data"
)

func TestItem_String(t *testing.T) {
	var tests = []struct {
		name     string
		expected string
		given    *data.Item
	}{
		{
			name:     "valid",
			expected: "Book: Programming Language Go($100)",
			given: &data.Item{
				ID:          "1",
				Name:        "Programming Language Go",
				Description: "Book of Go",
				Price:       100,
				Category:    data.Book,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.given.String()
			if actual != tt.expected {
				t.Errorf("(%s): expected %s, actual %s", tt.given, tt.expected, actual)
			}
		})
	}
}
