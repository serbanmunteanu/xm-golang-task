package main

import (
	"context"
	"time"

	"github.com/serbanmunteanu/xm-golang-task/config"
	http_framework "github.com/serbanmunteanu/xm-golang-task/http-framework"
	"github.com/serbanmunteanu/xm-golang-task/logger"
	"github.com/serbanmunteanu/xm-golang-task/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	logger.SetupErrorLog()
	cfg, err := config.LoadConfig("./config.yml")
	if err != nil {
		panic(err)
	}
	mongoClient, err := mongodb.NewMongoClient(cfg.MongoConfig)
	if err != nil {
		panic(err)
	}
	server := http_framework.NewHttpServer(cfg, mongoClient)
	server.Boot()
}

func initMongoDB(cfg config.MongoConfig) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Url))
	if err != nil {
		return nil, err
	}
	db := mongoClient.Database(cfg.Database)
	err = db.Client().Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	return db, nil
}
