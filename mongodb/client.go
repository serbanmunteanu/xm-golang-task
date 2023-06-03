package mongodb

import (
	"context"
	"time"

	"github.com/serbanmunteanu/xm-golang-task/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	database         *mongo.Database
	timeoutInSeconds time.Duration
}

func NewMongoClient(config config.MongoConfig) (*Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.Url))

	if err != nil {
		return nil, err
	}

	return &Client{
		client.Database(config.Database),
		time.Duration(config.TimeOutInSeconds),
	}, nil
}

func (client *Client) Insert(collection string, record interface{}) (interface{}, error) {
	coll := client.database.Collection(collection)

	ctx, cancel := context.WithTimeout(context.Background(), client.timeoutInSeconds*time.Second)
	defer cancel()

	result, err := coll.InsertOne(ctx, record)

	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

func (client *Client) FindOne(collection string, conditions bson.M, result interface{}) error {
	coll := client.database.Collection(collection)

	ctx, cancel := context.WithTimeout(context.Background(), client.timeoutInSeconds*time.Second)
	defer cancel()

	res := coll.FindOne(ctx, conditions).Decode(result)

	return res
}

func (client *Client) UpdateByID(collection string, id interface{}, model interface{}) error {
	coll := client.database.Collection(collection)

	ctx, cancel := context.WithTimeout(context.Background(), client.timeoutInSeconds*time.Second)
	defer cancel()

	_, err := coll.UpdateByID(
		ctx,
		id,
		model,
	)

	return err
}

func (client *Client) Find(collection string, filters bson.M) ([]bson.M, error) {
	coll := client.database.Collection(collection)

	ctx, cancel := context.WithTimeout(context.Background(), client.timeoutInSeconds*time.Second)
	defer cancel()

	cursor, err := coll.Find(ctx, filters)

	if err != nil {
		defer cursor.Close(ctx)
	}

	var results []bson.M
	for cursor.Next(ctx) {
		var result bson.M
		err = cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, err
}
