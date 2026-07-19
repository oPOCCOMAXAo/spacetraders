package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Service) Health(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
