package main

import (
	companyRepo "github.com/serbanmunteanu/xm-golang-task/company"
	"github.com/serbanmunteanu/xm-golang-task/config"
	"github.com/serbanmunteanu/xm-golang-task/di"
	httpframework "github.com/serbanmunteanu/xm-golang-task/http-framework"
	"github.com/serbanmunteanu/xm-golang-task/jwt"
	"github.com/serbanmunteanu/xm-golang-task/logger"
	"github.com/serbanmunteanu/xm-golang-task/mongodb"
	userRepo "github.com/serbanmunteanu/xm-golang-task/user"
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

	userRepository := userRepo.NewMongoRepository(mongoClient, cfg.MongoConfig.Collections.UserCollection)
	companyRepository := companyRepo.NewMongoRepository(mongoClient, cfg.MongoConfig.Collections.CompanyCollection)

	jwt := jwt.NewJwt(cfg.JwtConfig)

	diContainer := &di.DI{
		UserRepository:    userRepository,
		CompanyRepository: companyRepository,
		Jwt:               jwt,
	}

	server := httpframework.NewServer(cfg, diContainer)
	server.Boot()
}
