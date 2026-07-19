package strategy

import "github.com/opoccomaxao/spacetraders/pkg/services/spacetraders"

type Config struct{}

type Package struct {
	Service *Service
}

func MakePackage(
	config Config,
	spacetraders *spacetraders.Package,
) *Package {
	var res Package

	res.Service = NewService(
		spacetraders.Client,
	)

	return &res
}
