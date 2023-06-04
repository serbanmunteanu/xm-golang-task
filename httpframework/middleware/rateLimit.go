package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/serbanmunteanu/xm-golang-task/config"
	log "github.com/sirupsen/logrus"
)

type RateLimit struct {
}

func (rt *RateLimit) Register(router *gin.Engine, config *config.WebServerConfig) {
	concurrentRequests := config.ConcurrentRequests

	if concurrentRequests < 2 {
		concurrentRequests = 2
	}

	buffer := make(chan bool, concurrentRequests-1)

	router.Use(func(context *gin.Context) {
		select {
		case buffer <- true:
			context.Next()
			<-buffer
		default:
			log.Warning("Too many concurrent requests. Aborting.")
			context.Abort()
			context.JSON(http.StatusTooManyRequests, "test")
		}
	})
}
