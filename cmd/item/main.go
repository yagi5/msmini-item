package main

import (
	"context"
	"fmt"
	golog "log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/yagi5/msmini-item/di"
	"github.com/yagi5/msmini-item/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/reflection"
)

const (
	exitOK int = iota
	exitError
)

func main() {
	os.Exit(realMain(os.Args))
}

func realMain(args []string) int {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grpcLogger, err := log.New("ERROR")
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to setup logger: %s\n", err)
		return exitError
	}
	grpc_zap.ReplaceGrpcLogger(grpcLogger)

	container, err := di.New().GetContainer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to setup di container: %s\n", err)
		return exitError
	}

	reflection.Register(container.GRPCServer())

	wg, ctx := errgroup.WithContext(ctx)
	wg.Go(func() error { return container.GRPCServer().Serve(container.GRPCListener()) })

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-gracefulStop:
		golog.Println("received SIGTERM, exiting server gracefully")
	case <-ctx.Done():
	}

	container.GRPCServer().GracefulStop()
	cancel()
	if err := wg.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Unhandled error received: %s\n", err)
		return exitError
	}

	return exitOK
}
