package di

import (
	"net"
	"os"

	"github.com/yagi5/msmini-item/domain/repository"
	"github.com/yagi5/msmini-item/grpc"
	"github.com/yagi5/msmini-item/infrastructure/item"
	"github.com/yagi5/msmini-item/log"
	"github.com/yagi5/msmini-item/usecase"
	ggrpc "google.golang.org/grpc"
)

// Di returns container
type Di interface {
	Container() (*Container, error)
}

// Container is di container
type Container struct {
	grpcServer   *ggrpc.Server
	grpcListener net.Listener
	logger       *log.Logger
}

type client struct{}

// New returns di if
func New() Di {
	return &client{}
}

// Container is implementation of di interface
func (c *client) Container() (*Container, error) {
	repo := c.repository()
	logger, err := c.logger()
	if err != nil {
		return nil, err
	}
	usecase := c.usecase(logger, repo)
	grpcServer := c.grpcServer(usecase)
	grpcListener, err := c.grpcListener()
	if err != nil {
		return nil, err
	}

	return &Container{
		grpcServer:   grpcServer,
		grpcListener: grpcListener,
		logger:       logger,
	}, nil
}

func (c *client) repository() (repo repository.Item) {
	repo = item.NewDummyClient("/bin/items.csv")
	if os.Getenv("USE_CLOUDSQL") == "1" {
		repo = item.NewCloudSQLClient(nil)
	} else if os.Getenv("USE_SPANNER") == "1" {
		repo = item.NewSpannerClient(nil)
	}
	return
}

func (c *client) logger() (*log.Logger, error) {
	return log.New()
}

func (c *client) usecase(logger *log.Logger, repo repository.Item) usecase.Item {
	return usecase.New(logger, repo)
}

func (c *client) grpcServer(usecase usecase.Item) *ggrpc.Server {
	return grpc.New(usecase)
}

func (c *client) grpcListener() (net.Listener, error) {
	return net.Listen("tcp", ":10001")
}

// GRPCServer returns grpc.Server
func (c *Container) GRPCServer() *ggrpc.Server {
	return c.grpcServer
}

// GRPCListener returns net listener
func (c *Container) GRPCListener() net.Listener {
	return c.grpcListener
}

// Logger returns net logger
func (c *Container) Logger() *log.Logger {
	return c.logger
}
