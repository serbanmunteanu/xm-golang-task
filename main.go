package main

import (
	"github.com/serbanmunteanu/xm-golang-task/config"
	"github.com/serbanmunteanu/xm-golang-task/di"
	httpframework "github.com/serbanmunteanu/xm-golang-task/http-framework"
	"github.com/serbanmunteanu/xm-golang-task/jwt"
	"github.com/serbanmunteanu/xm-golang-task/logger"
	"github.com/serbanmunteanu/xm-golang-task/mongodb"
	"github.com/serbanmunteanu/xm-golang-task/user/repository"
)

func main() {
	logger.SetupErrorLog()
	cfg, err := config.LoadConfig("./config.yml")
	if err != nil {
		panic(err)
	}

	//@TODO: add switch for easy change to postgres repo type
	mongoClient, err := mongodb.NewMongoClient(cfg.MongoConfig)
	if err != nil {
		panic(err)
	}
	userRepository := repository.NewMongoUserRepository(mongoClient, cfg.MongoConfig.Collections.UserCollection)
	jwt := jwt.NewJwt(cfg.JwtConfig)

	diContainer := &di.DI{
		UserRepository: userRepository,
		Jwt:            jwt,
	}

	server := httpframework.NewServer(cfg, diContainer)
	server.Boot()
}
