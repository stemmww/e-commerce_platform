package usecase

import (
	"inventory/internal/model"
	"inventory/internal/repository"
)

type CategoryUsecase interface {
	GetAll() ([]model.Category, error)
	Create(model.Category) error
	Update(id int, category model.Category) error
	Delete(id int) error
	GetByID(id int) (model.Category, error)
}

type categoryUsecase struct {
	repo repository.CategoryRepository
}

func NewCategoryUsecase(r repository.CategoryRepository) CategoryUsecase {
	return &categoryUsecase{r}
}

func (u *categoryUsecase) GetAll() ([]model.Category, error) {
	return u.repo.GetAll()
}

func (u *categoryUsecase) Create(c model.Category) error {
	return u.repo.Create(c)
}

func (u *categoryUsecase) Update(id int, c model.Category) error {
	return u.repo.Update(id, c)
}

func (u *categoryUsecase) Delete(id int) error {
	return u.repo.Delete(id)
}

func (u *categoryUsecase) GetByID(id int) (model.Category, error) {
	return u.repo.GetByID(id)
}
