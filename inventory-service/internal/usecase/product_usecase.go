package usecase

import (
	"inventory/internal/model"
	"inventory/internal/repository"
)

type ProductUsecase interface {
	Create(product model.Product) error
	GetByID(id int) (*model.Product, error)
	GetAll() ([]model.Product, error)
	Update(product model.Product) error
	Delete(id int) error
}

type productUsecase struct {
	repo repository.ProductRepository
}

func NewProductUsecase(r repository.ProductRepository) ProductUsecase {
	return &productUsecase{r}
}

func (u *productUsecase) Create(product model.Product) error {
	return u.repo.Create(product)
}

func (u *productUsecase) GetByID(id int) (*model.Product, error) {
	return u.repo.GetByID(id)
}

func (u *productUsecase) GetAll() ([]model.Product, error) {
	return u.repo.GetAll()
}

func (u *productUsecase) Update(product model.Product) error {
	return u.repo.Update(product)
}

func (u *productUsecase) Delete(id int) error {
	return u.repo.Delete(id)
}
