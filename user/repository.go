package user

import (
	"github.com/serbanmunteanu/xm-golang-task/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository struct {
	client     *mongodb.Client
	collection string
}

func NewUserRepository(client *mongodb.Client, collection string) *UserRepository {
	return &UserRepository{
		client:     client,
		collection: collection,
	}
}

func (us *UserRepository) Insert(user *UserDbModel) error {
	id, err := us.client.Insert(us.collection, user)

	if err != nil {
		return err
	}

	user.ID = id.(primitive.ObjectID)

	return nil
}

func (us *UserRepository) FindOneBy(filters bson.M) (*UserDbModel, error) {
	var userModel UserDbModel

	err := us.client.FindOne(us.collection, filters, &userModel)

	if err != nil {
		return nil, err
	}

	return &userModel, nil
}

func (us *UserRepository) Find(filters bson.M) ([]UserDbModel, error) {
	var users []UserDbModel

	results, err := us.client.Find(us.collection, filters)

	if err != nil {
		return nil, err
	}

	for _, result := range results {
		var user UserDbModel
		bsonBytes, _ := bson.Marshal(result)
		bson.Unmarshal(bsonBytes, &user)
		users = append(users, user)
	}

	return users, nil
}
