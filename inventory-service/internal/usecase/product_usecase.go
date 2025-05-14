package usecase

import (
	"inventory/internal/model"
	"inventory/internal/repository"
)

type ProductUsecase interface {
	Create(product model.Product) error
	GetByID(id string) (*model.Product, error)
	GetAll() ([]model.Product, error)
	Update(product model.Product) error
	Delete(id string) error
}

// Define the struct (lowercase)
type productUsecase struct {
	repo repository.ProductRepository
}

func NewProductUsecase(r repository.ProductRepository) ProductUsecase {
	return &productUsecase{repo: r}
}

func (u *productUsecase) Create(product model.Product) error {
	return u.repo.Create(product)
}

func (u *productUsecase) GetByID(id string) (*model.Product, error) {
	return u.repo.GetByID(id)
}

func (u *productUsecase) GetAll() ([]model.Product, error) {
	return u.repo.GetAll()
}

func (u *productUsecase) Update(product model.Product) error {
	return u.repo.Update(product)
}

func (u *productUsecase) Delete(id string) error {
	return u.repo.Delete(id)
}
