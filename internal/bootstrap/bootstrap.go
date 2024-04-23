package bootstrap

import (
	"errors"
	"fmt"
	"github.com/samber/do"
	"go.uber.org/zap"
	"ltp-api/internal/config"
	"ltp-api/internal/services"
	"ltp-api/internal/transport/http/handlers"
	"ltp-api/internal/transport/http/server"
	"ltp-api/pkg/kraken"
	"ltp-api/pkg/redis"
)

func Bootstrap(injector *do.Injector) {

	do.Provide[*zap.Logger](injector, func(injector *do.Injector) (*zap.Logger, error) {
		logger, err := zap.NewDevelopment()

		if err != nil {
			return nil, errors.New(fmt.Sprintf("can't initialize zap logger: %v", err))
		}

		zap.ReplaceGlobals(logger)
		return logger, nil
	})

	do.Provide[*config.Config](injector, func(i *do.Injector) (*config.Config, error) {
		return config.New(".", "config.yaml")
	})

	do.Provide[*redis.Cache](injector, func(i *do.Injector) (*redis.Cache, error) {
		cfg := do.MustInvoke[*config.Config](i)

		return redis.New(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password), nil
	})

	do.Provide[kraken.Client](injector, func(i *do.Injector) (kraken.Client, error) {
		cfg := do.MustInvoke[*config.Config](i)

		return kraken.New(&cfg.Kraken), nil
	})

	// Provide All Repositories here

	// Provide All Services here
	do.Provide[services.KrakenService](injector, services.NewKrakenService)

	//Register Public Handlers here
	do.ProvideNamed[[]handlers.Handler](injector, "publicHandlers", func(i *do.Injector) ([]handlers.Handler, error) {
		return []handlers.Handler{
			handlers.NewLTPHandler(i),
		}, nil
	})

	//Register Authentication Handlers here

	do.Provide[*server.Server](injector, server.New)

}
