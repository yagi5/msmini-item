package di

import (
	"net"
	"os"

	"github.com/yagi5/msmini-item/grpc"
	"github.com/yagi5/msmini-item/infrastructure/item"
	"github.com/yagi5/msmini-item/log"
	"github.com/yagi5/msmini-item/usecase"
	"go.uber.org/zap"
	ggrpc "google.golang.org/grpc"
)

// Di returns container
type Di interface {
	GetContainer() (*Container, error)
}

// Container is di container
type Container struct {
	grpcServer   *ggrpc.Server
	grpcListener net.Listener
	Logger       *zap.Logger
}

type client struct{}

// New returns di if
func New() Di {
	return &client{}
}

// GetContainer is implementation of di interface
func (c *client) GetContainer() (*Container, error) {
	repo := item.NewDummyClient()
	if os.Getenv("USE_CLOUDSQL") == "1" {
		repo = item.NewCloudSQLClient(nil)
	} else if os.Getenv("USE_SPANNER") == "1" {
		repo = item.NewSpannerClient(nil)
	}

	logger, err := log.New()
	if err != nil {
		return nil, err
	}

	usecase := usecase.New(logger, repo)

	grpcS := grpc.New(usecase)
	grpcL, err := net.Listen("tcp", ":10001")
	if err != nil {
		return nil, err
	}

	return &Container{
		grpcServer:   grpcS,
		grpcListener: grpcL,
	}, nil
}

// GRPCServer returns grpc.Server
func (c *Container) GRPCServer() *ggrpc.Server {
	return c.grpcServer
}

// GRPCListener returns net listener
func (c *Container) GRPCListener() net.Listener {
	return c.grpcListener
}
