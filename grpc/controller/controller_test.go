package controller

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/go-cmp/cmp"
	"github.com/yagi5/msmini-item/domain/data"
	"github.com/yagi5/msmini-item/usecase"
	pb "github.com/yagi5/msmini-proto/client/proto/item"
)

func TestNew(t *testing.T) {
	te := time.Date(2019, 1, 2, 12, 0, 30, 0, time.Local)
	tw, err := ptypes.TimestampProto(te)
	if err != nil {
		t.Fatal(err)
	}
	var tests = []struct {
		name       string
		input      *pb.SearchByNameRequest
		searchMock func(context.Context, *usecase.ItemSearchInput) ([]*data.Item, error)
		wantErr    bool
		want       *pb.SearchByNameResponse
	}{
		{
			name:  "return error",
			input: &pb.SearchByNameRequest{Name: ""},
			searchMock: func(context.Context, *usecase.ItemSearchInput) ([]*data.Item, error) {
				return nil, errors.New("")
			},
			wantErr: true,
			want:    nil,
		},
		{
			name:  "succeeds",
			input: &pb.SearchByNameRequest{Name: ""},
			searchMock: func(context.Context, *usecase.ItemSearchInput) ([]*data.Item, error) {
				return []*data.Item{
					{
						ID:          "1",
						Name:        "name",
						Description: "description",
						Price:       100,
						Category:    data.Category("Book"),
						CreatedAt:   te,
						UpdatedAt:   te,
					},
				}, nil
			},
			want: &pb.SearchByNameResponse{Items: []*pb.Item{
				{
					Id:          "1",
					Name:        "name",
					Description: "description",
					Price:       100,
					Category:    "Book",
					CreatedAt:   tw,
					UpdatedAt:   tw,
				},
			}},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			i := New(&usecase.ItemMock{SearchMock: tt.searchMock})
			items, err := i.SearchByName(context.TODO(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("failed wantErr: %v, got %v", tt.wantErr, err)
			}
			if diff := cmp.Diff(items, tt.want); diff != "" {
				t.Errorf("failed %s", diff)
			}
		})
	}
}
