package app

import (
	"context"
	"os"
	"os/signal"

	"github.com/opoccomaxao/gopkg/pkg/services/ginserver"
	"github.com/opoccomaxao/gopkg/pkg/services/lifecycle"
	"github.com/opoccomaxao/gopkg/pkg/services/logger"
	"github.com/opoccomaxao/spacetraders/pkg/config"
	"github.com/opoccomaxao/spacetraders/pkg/endpoints"
	"github.com/opoccomaxao/spacetraders/pkg/services/spacetraders"
	"github.com/opoccomaxao/spacetraders/pkg/services/strategy"
)

func Run() error {
	appCtx, appCancelCause := context.WithCancelCause(context.Background())
	defer appCancelCause(nil)

	appCtx, cancel := signal.NotifyContext(appCtx, os.Interrupt)
	defer cancel()

	config, err := config.Load()
	if err != nil {
		return err
	}

	logger := logger.MakePackage(config.Logger)

	lifecycle := lifecycle.MakePackage(
		config.Lifecycle,
		appCancelCause,
		logger,
	)

	server := ginserver.MakePackage(
		config.Server,
		lifecycle,
		logger,
	)

	spacetraders, err := spacetraders.MakePackage(
		config.SpaceTraders,
		lifecycle,
		logger,
	)
	if err != nil {
		logger.Error(err)

		return err
	}

	strategy := strategy.MakePackage(
		config.Strategy,
		spacetraders,
	)

	err = endpoints.MakePackage(
		config.Endpoints,
		server,
		strategy,
	)
	if err != nil {
		logger.Error(err)

		return err
	}

	err = lifecycle.Service.Serve(appCtx, appCancelCause)
	if err != nil {
		logger.Error(err)

		return err
	}

	<-appCtx.Done()

	err = lifecycle.Service.Shutdown(appCtx)
	if err != nil {
		logger.Error(err)

		return err
	}

	return nil
}
