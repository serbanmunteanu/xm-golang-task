package user

import (
	"github.com/serbanmunteanu/xm-golang-task/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mongoUserRepository struct {
	client     *mongodb.Client
	collection string
}

func (us *mongoUserRepository) Create(user *Model) error {
	id, err := us.client.Insert(us.collection, user)

	if err != nil {
		return err
	}

	user.ID = id.(primitive.ObjectID)

	return nil
}

func (us *mongoUserRepository) Read(email string) (*Model, error) {
	var userModel Model

	err := us.client.FindOne(us.collection, bson.M{"email": email}, &userModel)

	if err != nil {
		return nil, err
	}

	return &userModel, nil
}

func NewMongoRepository(client *mongodb.Client, collection string) IRepository {
	return &mongoUserRepository{
		client:     client,
		collection: collection,
	}
}
