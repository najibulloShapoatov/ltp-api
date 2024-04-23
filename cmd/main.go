package main

import (
	"context"
	"flag"
	"github.com/samber/do"
	"go.uber.org/zap"
	"ltp-api/internal/bootstrap"
	"ltp-api/internal/transport/http/server"
	"ltp-api/internal/utils"
	"sync"
	"time"
)

func main() {
	flag.Parse()

	now := time.Now()
	injector := do.New()
	ctx := context.Background()
	wg := &sync.WaitGroup{}

	do.Provide[*sync.WaitGroup](injector, func(i *do.Injector) (*sync.WaitGroup, error) {
		return wg, nil
	})

	do.Provide[context.Context](injector, func(i *do.Injector) (context.Context, error) {
		return ctx, nil
	})

	bootstrap.Bootstrap(injector)

	logger := do.MustInvoke[*zap.Logger](injector)

	logger.Sugar().Info("Starting application...")

	s := do.MustInvoke[*server.Server](injector)

	go s.Run()

	zap.S().Infof("Up and running (%s)", time.Since(now))
	zap.S().Infof("Got %s signal. Shutting down...", <-utils.WaitTermSignal())

	if err := s.Shutdown(ctx); err != nil {
		zap.S().Errorf("Error stopping server: %s", err)
	}

	if err := injector.Shutdown(); err != nil {
		zap.S().Errorf("Error stopping injector: %s", err)
	}

	wg.Wait()
}
