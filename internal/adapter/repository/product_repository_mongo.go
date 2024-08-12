package repository

import (
	"context"
	"go-hexagon/internal/core/domain/entity"
	"go-hexagon/internal/core/port"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepositoryMongo struct {
    DB *mongo.Collection
}

func NewProductRepositoryMongo(db *mongo.Database) port.ProductRepository {
    return &ProductRepositoryMongo{DB: db.Collection("products")}
}

func (r *ProductRepositoryMongo) Create(product *entity.Product) error {
    _, err := r.DB.InsertOne(context.Background(), product)
    return err
}

func (r *ProductRepositoryMongo) Update(product *entity.Product) error {
    filter := bson.M{"_id": product.ID}
    update := bson.M{"$set": product}
    _, err := r.DB.UpdateOne(context.Background(), filter, update)
    return err
}

func (r *ProductRepositoryMongo) GetByID(id uint) (*entity.Product, error) {
    var product entity.Product
    filter := bson.M{"_id": id}
    err := r.DB.FindOne(context.Background(), filter).Decode(&product)
    return &product, err
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

func (r *ProductRepositoryMongo) Delete(id uint) error {
    filter := bson.M{"_id": id}
    _, err := r.DB.DeleteOne(context.Background(), filter)
    return err
}
