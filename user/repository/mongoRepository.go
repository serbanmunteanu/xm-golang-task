package repository

import (
	"github.com/serbanmunteanu/xm-golang-task/mongodb"
	"github.com/serbanmunteanu/xm-golang-task/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mongoUserRepository struct {
	client     *mongodb.Client
	collection string
}

func NewMongoUserRepository(client *mongodb.Client, collection string) UserRepository {
	return &mongoUserRepository{
		client:     client,
		collection: collection,
	}
}

func (us *mongoUserRepository) Create(user *user.UserDbModel) error {
	id, err := us.client.Insert(us.collection, user)

	if err != nil {
		return err
	}

	user.ID = id.(primitive.ObjectID)

	return nil
}

func (us *mongoUserRepository) Read(email string) (*user.UserDbModel, error) {
	var userModel user.UserDbModel

	err := us.client.FindOne(us.collection, bson.M{"email": email}, &userModel)

	if err != nil {
		return nil, err
	}

	return &userModel, nil
}
