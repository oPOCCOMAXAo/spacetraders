package strategy

import "github.com/opoccomaxao/spacetraders/pkg/services/spacetraders"

type Service struct {
	spacetraders *spacetraders.Client
}

func NewService(
	spacetraders *spacetraders.Client,
) *Service {
	return &Service{
		spacetraders: spacetraders,
	}
}
