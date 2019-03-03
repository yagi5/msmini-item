package grpc

import (
	"github.com/yagi5/msmini-item/grpc/controller"
	"github.com/yagi5/msmini-item/usecase"
	pb "github.com/yagi5/msmini-proto/client/proto/item"
	"google.golang.org/grpc"
)

// New returns grpc server
func New(usecase usecase.Item) *grpc.Server {
	s := grpc.NewServer()
	item := controller.New(usecase)
	pb.RegisterItemServiceServer(s, item)
	return s
}
