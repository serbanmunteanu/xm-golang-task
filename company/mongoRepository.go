package company

import (
	"github.com/serbanmunteanu/xm-golang-task/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mongoCompanyRepository struct {
	client     *mongodb.Client
	collection string
}

func (m mongoCompanyRepository) Create(company *Model) error {
	id, err := m.client.Insert(m.collection, company)
	if err != nil {
		return err
	}
	company.ID = id.(primitive.ObjectID)
	return nil
}

func (m mongoCompanyRepository) Read(id primitive.ObjectID) (*Model, error) {
	var model Model
	err := m.client.FindOne(m.collection, bson.M{"_id": id}, &model)
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (m mongoCompanyRepository) Count(name string) (int64, error) {
	return m.client.CountDocuments(m.collection, bson.M{"name": name})
}

func (m mongoCompanyRepository) Update(id primitive.ObjectID, updates interface{}) error {
	return m.client.UpdateByID(m.collection, id, bson.M{"$set": updates})
}

func (m mongoCompanyRepository) Delete(id primitive.ObjectID) error {
	return m.client.Delete(m.collection, id)
}

func NewMongoRepository(client *mongodb.Client, collection string) CompanyRepository {
	return &mongoCompanyRepository{
		client:     client,
		collection: collection,
	}
}
