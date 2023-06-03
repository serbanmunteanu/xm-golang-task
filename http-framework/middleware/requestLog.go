package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/serbanmunteanu/xm-golang-task/config"
	"github.com/serbanmunteanu/xm-golang-task/logger"
	log "github.com/sirupsen/logrus"
)

type RequestLog struct {
}

func (rl *RequestLog) Register(router *gin.Engine, config *config.WebServerConfig) {
	accessLog := logger.SetupAccessLog()

	router.Use(func(context *gin.Context) {
		context.Next()

		accessLog.WithFields(log.Fields{
			"ip":     context.Request.RemoteAddr,
			"url":    context.Request.URL.Path,
			"method": context.Request.Method,
			"status": context.Writer.Status(),
		}).Info("")
	})
}
