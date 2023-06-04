package grpc

import (
	"github.com/serbanmunteanu/xm-golang-task/config"
	"github.com/serbanmunteanu/xm-golang-task/di"
	"github.com/serbanmunteanu/xm-golang-task/server"
)

type GrpcServer struct {
	config *config.WebServerConfig
	di     *di.DI
}

func (g *GrpcServer) Boot() {
	//@TODO: change framework
}

func NewServer(config *config.WebServerConfig, di *di.DI) server.IServer {
	return &GrpcServer{
		config: config,
		di:     di,
	}
}
