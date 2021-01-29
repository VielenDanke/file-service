package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/vielendanke/file-service/internal/app/fileservice"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.DefaultLogger = logger.NewLogger(logger.WithLevel(logger.TraceLevel))

	errCh := make(chan error, 1)

	go func() {
		sigCh := make(chan os.Signal)
		signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
		errCh <- fmt.Errorf("%s", <-sigCh)
	}()

	go fileservice.StartFileService(ctx, errCh)

	logger.Infof(ctx, "Service terminated: %v", <-errCh)
}
