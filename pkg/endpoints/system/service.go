package system

import (
	"github.com/gin-gonic/gin"
	"github.com/opoccomaxao/spacetraders/pkg/endpoints/internal"
)

type Service struct{}

func New() internal.EndpointsService {
	return &Service{}
}

func (s *Service) Register(router gin.IRouter) error {
	router.GET("/api/health", s.Health)

	return nil
}
