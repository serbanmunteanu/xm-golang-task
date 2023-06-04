package main

import (
	"context"
	"fmt"

	companyRepo "github.com/serbanmunteanu/xm-golang-task/company"
	"github.com/serbanmunteanu/xm-golang-task/config"
	"github.com/serbanmunteanu/xm-golang-task/di"
	httpframework "github.com/serbanmunteanu/xm-golang-task/httpframework"
	"github.com/serbanmunteanu/xm-golang-task/jwt"
	"github.com/serbanmunteanu/xm-golang-task/kafka"
	"github.com/serbanmunteanu/xm-golang-task/logger"
	"github.com/serbanmunteanu/xm-golang-task/mongodb"
	userRepo "github.com/serbanmunteanu/xm-golang-task/user"
	log "github.com/sirupsen/logrus"
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
	producer := kafka.NewProducer(cfg.KafkaConfig)
	consumer := kafka.NewConsumer(cfg.KafkaConfig)

	go func() {
		for {
			m, err := consumer.Reader.ReadMessage(context.Background())
			if err != nil {
				log.Fatal("failed to read message: ", err)
			}
			fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		}
	}()

	userRepository := userRepo.NewMongoRepository(mongoClient, cfg.MongoConfig.Collections.UserCollection)
	companyRepository := companyRepo.NewMongoRepository(mongoClient, cfg.MongoConfig.Collections.CompanyCollection)

	jwt := jwt.NewJwt(cfg.JwtConfig)

	diContainer := &di.DI{
		UserRepository:    userRepository,
		CompanyRepository: companyRepository,
		Jwt:               jwt,
		Producer:          producer,
	}

	server := httpframework.NewServer(cfg, diContainer)
	//server = grpc.NewServer(cfg, diContainer) // switch to grpc easily
	server.Boot()
}
