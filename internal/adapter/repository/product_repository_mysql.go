package repository

import (
	"fmt"
	"go-hexagon/internal/core/domain/entity"
	"go-hexagon/internal/core/port"

	"gorm.io/gorm"
)

type ProductRepositoryMySQL struct {
	DB *gorm.DB
}

func NewProductRepositoryMySQL(db *gorm.DB) port.ProductRepository {
	return &ProductRepositoryMySQL{DB: db}
}

func (r *ProductRepositoryMySQL) Create(product *entity.Product) error {
	return r.DB.Table("products").Create(&product).Error
}

func (r *ProductRepositoryMySQL) Update(product *entity.Product) error {
	var existingProduct entity.Product
	if err := r.DB.Table("products").First(&existingProduct, product.MySQLID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("ID Not Found")
		}
		return err
	}

	return r.DB.Table("products").Save(product).Error
}

func (r *ProductRepositoryMySQL) GetByID(id interface{}) (*entity.Product, error) {
	var product entity.Product

	idUint, ok := id.(uint)
	if !ok {
		return nil, fmt.Errorf("Invalid ID type for MySQL")
	}

	err := r.DB.Table("products").First(&product, idUint).Error
	return &product, err
}

func (r *ProductRepositoryMySQL) List() ([]entity.Product, error) {
	var products []entity.Product
	err := r.DB.Table("products").Find(&products).Error
	return products, err
}

func (r *ProductRepositoryMySQL) Delete(id interface{}) error {
	idUint, ok := id.(uint)
	if !ok {
		return fmt.Errorf("Invalid ID type for MySQL")
	}

	var product entity.Product
	if err := r.DB.Table("products").First(&product, idUint).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("ID Not Found")
		}
		return err
	}

	return r.DB.Table("products").Delete(&product).Error
}
