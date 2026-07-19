package endpoints

import (
	"github.com/opoccomaxao/gopkg/pkg/services/ginserver"
	"github.com/opoccomaxao/spacetraders/pkg/endpoints/internal"
	"github.com/opoccomaxao/spacetraders/pkg/endpoints/system"
	"github.com/opoccomaxao/spacetraders/pkg/services/strategy"
)

type Config struct{}

func MakePackage(
	config Config,
	server *ginserver.Package,
	strategy *strategy.Package,
) error {
	regs := []internal.EndpointsService{
		system.New(),
	}

	for _, reg := range regs {
		err := reg.Register(server.Router)
		if err != nil {
			return err
		}
	}

	return nil
}
