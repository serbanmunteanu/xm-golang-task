package company

import "go.mongodb.org/mongo-driver/bson/primitive"

type IRepository interface {
	Create(company *Model) error
	Read(id primitive.ObjectID) (*Model, error)
	Count(name string) (int64, error)
	Update(id primitive.ObjectID, updates interface{}) error
	Delete(id primitive.ObjectID) error
}
