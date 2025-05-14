package usecase

import (
	"user-service/internal/model"
	"user-service/internal/repository"
)

type UserUsecase interface {
	Register(user *model.User) error
	Authenticate(email, password string) (*model.User, error)
	GetByID(id string) (*model.User, error)
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(r repository.UserRepository) UserUsecase {
	return &userUsecase{repo: r}
}

func (u *userUsecase) Register(user *model.User) error {
	return u.repo.Register(user)
}

func (u *userUsecase) Authenticate(email, password string) (*model.User, error) {
	return u.repo.Authenticate(email, password)
}

func (u *userUsecase) GetByID(id string) (*model.User, error) {
	return u.repo.GetByID(id)
}
