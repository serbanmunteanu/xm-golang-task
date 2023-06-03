package http_framework

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/serbanmunteanu/xm-golang-task/config"
	"github.com/serbanmunteanu/xm-golang-task/di"
	"github.com/serbanmunteanu/xm-golang-task/http-framework/middleware"
	log "github.com/sirupsen/logrus"
)

type HttpServer struct {
	config *config.WebServerConfig
	di     *di.DI
}

func NewServer(config *config.WebServerConfig, di *di.DI) *HttpServer {
	return &HttpServer{
		config: config,
		di:     di,
	}
}

func (hs *HttpServer) Boot() {
	router := gin.New()

	middlewares := []middleware.RouterMiddleware{
		&middleware.RateLimit{},
		&middleware.RequestLog{},
	}

	for _, middle := range middlewares {
		middle.Register(router, hs.config)
	}

	router.Use(gin.Recovery())

	Initialize(router, hs.di)

	log.Info("Starting server on port ", hs.config.ServerPort)

	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", hs.config.ServerPort),
		Handler: router,
	}

	go gracefulShutdown(
		server,
		quit,
		done,
	)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	log.Info("Graceful shutdown completed")
}

func gracefulShutdown(server *http.Server, quit <-chan os.Signal, done chan<- bool) {
	<-quit
	log.Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %v", err)
	}

	close(done)
}
