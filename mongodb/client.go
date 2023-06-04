package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/serbanmunteanu/xm-golang-task/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (client *Client) CountDocuments(collection string, filters bson.M) (int64, error) {
	coll := client.database.Collection(collection)

	ctx, cancel := context.WithTimeout(context.Background(), client.timeoutInSeconds*time.Second)
	defer cancel()

	count, err := coll.CountDocuments(ctx, filters)

	return count, err
}

func (client *Client) Delete(collection string, id primitive.ObjectID) error {
	coll := client.database.Collection(collection)

	ctx, cancel := context.WithTimeout(context.Background(), client.timeoutInSeconds*time.Second)
	defer cancel()

	result, err := coll.DeleteOne(ctx, bson.M{"_id": id})

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no documents were deleted")
	}

	return nil
}
