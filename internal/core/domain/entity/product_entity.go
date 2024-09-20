package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	MySQLID uint               `json:"id,omitempty" gorm:"primaryKey;autoIncrement;column:id" bson:"-"`
	MongoID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty" gorm:"-"`
	Name    string             `json:"name" gorm:"column:name" bson:"name"`
	Stock   int                `json:"stock" gorm:"column:stock" bson:"stock"`
}
