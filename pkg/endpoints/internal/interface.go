package internal

import "github.com/gin-gonic/gin"

type EndpointsService interface {
	Register(gin.IRouter) error
}
