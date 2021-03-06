package main

import (
	"context"
	"fmt"
	golog "log"
	"os"
	"os/signal"
	"syscall"
	"time"

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

var version string

func main() {
	os.Exit(realMain(os.Args))
}

func realMain(args []string) int {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		loc = time.FixedZone("Asia/Tokyo", 9*60*60)
	}
	time.Local = loc

	logger, err := log.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to setup logger: %s\n", err)
		return exitError
	}
	grpc_zap.ReplaceGrpcLogger(logger.L)

	container, err := di.New().Container()
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
