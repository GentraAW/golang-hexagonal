package port

import "go-hexagon/internal/core/domain/entity"

type ProductRepository interface {
	Create(product *entity.Product) error
	Update(product *entity.Product) error
	GetByID(id uint) (*entity.Product, error)
	List() ([]entity.Product, error)
	Delete(id uint) error
}
