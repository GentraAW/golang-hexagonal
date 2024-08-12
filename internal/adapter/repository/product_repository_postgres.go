package repository

import (
	"fmt"
	"go-hexagon/internal/core/domain/entity"
	"go-hexagon/internal/core/port"

	"gorm.io/gorm"
)

type ProductRepositoryPostgres struct {
	DB *gorm.DB
}

func NewProductRepositoryPostgres(db *gorm.DB) port.ProductRepository {
	return &ProductRepositoryPostgres{DB: db}
}

func (r *ProductRepositoryPostgres) Create(product *entity.Product) error {
	return r.DB.Create(product).Error
}

func (r *ProductRepositoryPostgres) Update(product *entity.Product) error {
    var existingProduct entity.Product

    // Cek apakah produk dengan ID yang diberikan ada di database
    if err := r.DB.First(&existingProduct, product.ID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return fmt.Errorf("ID: tidak tersedia")
        }
        return err
    }

    // Lakukan update jika produk ditemukan
    return r.DB.Save(product).Error
}

func (r *ProductRepositoryPostgres) GetByID(id uint) (*entity.Product, error) {
	var product entity.Product
	err := r.DB.First(&product, id).Error
	return &product, err
}

func (r *ProductRepositoryPostgres) List() ([]entity.Product, error) {
	var products []entity.Product
	err := r.DB.Find(&products).Error
	return products, err
}

func (r *ProductRepositoryPostgres) Delete(id uint) error {
    var product entity.Product

    // Cek apakah produk dengan ID yang diberikan ada di database
    if err := r.DB.First(&product, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return fmt.Errorf("ID: tidak tersedia")
        }
        return err
    }

    // Jika produk ditemukan, lakukan delete
    return r.DB.Delete(&product).Error
}