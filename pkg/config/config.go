package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/opoccomaxao/gopkg/pkg/services/ginserver"
	"github.com/opoccomaxao/gopkg/pkg/services/lifecycle"
	"github.com/opoccomaxao/gopkg/pkg/services/logger"
	"github.com/opoccomaxao/spacetraders/pkg/endpoints"
	"github.com/opoccomaxao/spacetraders/pkg/services/spacetraders"
	"github.com/opoccomaxao/spacetraders/pkg/services/strategy"
	"github.com/pkg/errors"
)

type Config struct {
	Logger       logger.Config       `envPrefix:"LOGGER_"`
	Lifecycle    lifecycle.Config    `envPrefix:"LIFECYCLE_"`
	Server       ginserver.Config    `envPrefix:"SERVER_"`
	SpaceTraders spacetraders.Config `envPrefix:"SPACETRADERS_"`
	Strategy     strategy.Config     `envPrefix:"STRATEGY_"`
	Endpoints    endpoints.Config    `envPrefix:"ENDPOINTS_"`
}

func Load() (*Config, error) {
	var res Config

	err := env.ParseWithOptions(&res, env.Options{
		RequiredIfNoDef:       false,
		UseFieldNameByDefault: false,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &res, nil
}
