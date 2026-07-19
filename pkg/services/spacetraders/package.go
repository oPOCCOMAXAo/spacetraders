package spacetraders

import (
	"log/slog"

	"github.com/opoccomaxao/gopkg/pkg/services/lifecycle"
	"github.com/opoccomaxao/gopkg/pkg/services/logger"
)

type Config struct {
	Token         string `env:"TOKEN,required"`
	Host          string `env:"HOST"           envDefault:"https://api.spacetraders.io/v2"`
	VerboseErrors bool   `env:"VERBOSE_ERRORS" envDefault:"false"`
}

type Package struct {
	Client *Client
}

func MakePackage(
	config Config,
	lifecycle *lifecycle.Package,
	logger *logger.Package,
) (*Package, error) {
	var (
		res Package
		err error
	)

	res.Client, err = NewClient(
		config,
		logger.Logger.With(slog.String("service", "spacetraders")),
	)
	if err != nil {
		return nil, err
	}

	lifecycle.Service.RegisterService(res.Client)

	return &res, nil
}
