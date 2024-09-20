package port

import "go-hexagon/internal/core/domain/entity"

type ProductRepository interface {
	Create(product *entity.Product) error
	Update(product *entity.Product) error
	GetByID(id interface{}) (*entity.Product, error)
	List() ([]entity.Product, error)
	Delete(id interface{}) error
}
