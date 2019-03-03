package controller

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"github.com/yagi5/msmini-item/usecase"
	pb "github.com/yagi5/msmini-proto/client/proto/item"
)

type item struct {
	usecase usecase.Item
}

// New returns itemcontroller
func New(usecase usecase.Item) pb.ItemServiceServer {
	return &item{usecase: usecase}
}

// SearchByName implements item interface and proto-defined function
func (i *item) SearchByName(ctx context.Context, req *pb.SearchByNameRequest) (*pb.SearchByNameResponse, error) {
	in := &usecase.ItemSearchInput{Name: req.Name}
	items, err := i.usecase.Search(ctx, in)
	if err != nil {
		return nil, errors.Wrap(err, "search failed")
	}

	var is []*pb.Item
	for _, item := range items {

		ca, err := ptypes.TimestampProto(item.CreatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "convert createdAt into timestampproto failed")
		}
		ua, err := ptypes.TimestampProto(item.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "convert updatedAt into timestampproto failed")
		}
		i := &pb.Item{
			Id:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
			Category:    string(item.Category),
			CreatedAt:   ca,
			UpdatedAt:   ua,
		}
		is = append(is, i)
	}
	return &pb.SearchByNameResponse{Items: is}, nil
}
