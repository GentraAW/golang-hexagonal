package service

import (
	"go-hexagon/internal/core/domain/entity"
	"go-hexagon/internal/core/port"
)

type ProductService struct {
	Repo port.ProductRepository
}

func NewProductService(repo port.ProductRepository) *ProductService {
	return &ProductService{Repo: repo}
}

func (s *ProductService) CreateProduct(product *entity.Product) error {
	return s.Repo.Create(product)
}

func (s *ProductService) UpdateProduct(product *entity.Product) error {
	return s.Repo.Update(product)
}

func (s *ProductService) GetProductByID(id interface{}) (*entity.Product, error) {
	return s.Repo.GetByID(id)
}

func (s *ProductService) ListProducts() ([]entity.Product, error) {
	return s.Repo.List()
}

func (s *ProductService) DeleteProduct(id interface{}) error {
	return s.Repo.Delete(id)
}
