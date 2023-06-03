package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/serbanmunteanu/xm-golang-task/config"
)

type RouterMiddleware interface {
	Register(router *gin.Engine, config *config.WebServerConfig)
}
