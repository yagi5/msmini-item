package usecase

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"github.com/yagi5/msmini-item/domain/data"
	"github.com/yagi5/msmini-item/domain/repository"
)

func TestSearch(t *testing.T) {
	var tests = []struct {
		name           string
		wantErr        bool
		mockSearchFunc func(context.Context, string, int64) ([]*data.Item, error)
		given          *ItemSearchInput
		expected       []*data.Item
	}{
		{
			name:    "name is empty",
			wantErr: true,
			given:   &ItemSearchInput{},
		},
		{
			name:  "Limit has default value",
			given: &ItemSearchInput{Name: "test"},
			mockSearchFunc: func(ctx context.Context, name string, limit int64) ([]*data.Item, error) {
				if limit != 20 {
					t.Fatalf("limit is not 20: %d", limit)
				}
				return []*data.Item{}, nil
			},
			expected: []*data.Item{},
		},
		{
			name:    "error returned",
			wantErr: true,
			given:   &ItemSearchInput{Name: "test"},
			mockSearchFunc: func(ctx context.Context, name string, limit int64) ([]*data.Item, error) {
				return nil, errors.New("test")
			},
		},
		{
			name:  "succeeded",
			given: &ItemSearchInput{Name: "test"},
			mockSearchFunc: func(ctx context.Context, name string, limit int64) ([]*data.Item, error) {
				return []*data.Item{&data.Item{Name: "test"}}, nil
			},
			expected: []*data.Item{&data.Item{Name: "test"}},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repos := Repositories{Item: &repository.ItemMock{SearchByNameMock: tt.mockSearchFunc}}
			iu := ItemUsecase{repos: repos}
			actual, err := iu.Search(context.Background(), tt.given)
			if (err != nil) != tt.wantErr {
				if err != nil {
					t.Fatal(err)
				}
			}
			if diff := cmp.Diff(actual, tt.expected); diff != "" {
				t.Errorf("failed: %s", diff)
			}
		})
	}
}
