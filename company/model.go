package company

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	ID                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name              string             `json:"name" bson:"name" binding:"required,max=15"`
	Description       string             `json:"description,omitempty" bson:"description,omitempty" binding:"max=3000"`
	AmountOfEmployees int                `json:"amountOfEmployees" bson:"amountOfEmployees" binding:"required"`
	Registered        bool               `json:"registered" bson:"registered" binding:"required"`
	Type              string             `json:"type" bson:"type" binding:"required,oneof=Corporations NonProfit Cooperative SoleProprietorship"`
	CreatedAt         time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt         time.Time          `json:"updatedAt" bson:"updatedAt"`
}
