package repository

import (
	"context"
	"fmt"
	"go-hexagon/internal/core/domain/entity"
	"go-hexagon/internal/core/port"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepositoryMongo struct {
	DB *mongo.Collection
}

func NewProductRepositoryMongo(db *mongo.Database) port.ProductRepository {
	return &ProductRepositoryMongo{DB: db.Collection("products")}
}

func (r *ProductRepositoryMongo) Create(product *entity.Product) error {
	product.MongoID = primitive.NewObjectID()
	_, err := r.DB.InsertOne(context.Background(), product)
	return err
}

func (r *ProductRepositoryMongo) Update(product *entity.Product) error {
	filter := bson.M{"_id": product.MongoID}
	update := bson.M{"$set": product}
	result, err := r.DB.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("ID Not Found")
	}
	return nil
}

func (r *ProductRepositoryMongo) GetByID(id interface{}) (*entity.Product, error) {
	var product entity.Product
	if idObj, ok := id.(primitive.ObjectID); ok {
		filter := bson.M{"_id": idObj}
		err := r.DB.FindOne(context.Background(), filter).Decode(&product)
		return &product, err
	}
	return nil, fmt.Errorf("Invalid ID type for MongoDB")
}

func (r *ProductRepositoryMongo) List() ([]entity.Product, error) {
	var products []entity.Product
	cursor, err := r.DB.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(context.Background(), &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepositoryMongo) Delete(id interface{}) error {
	objectID, ok := id.(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("Invalid ID type for MongoDB")
	}

	filter := bson.M{"_id": objectID}
	res := r.DB.FindOne(context.Background(), filter)
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return fmt.Errorf("ID Not Found")
		}
		return res.Err()
	}

	_, err := r.DB.DeleteOne(context.Background(), filter)
	return err
}
