package usecase

import (
	"user-service/internal/model"
	"user-service/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Register(user *model.User) error
	Login(email, password string) (model.User, error)
	GetProfile(id int) (model.User, error)
	UpdateUser(user model.User) error
	DeleteUser(id int) error
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) Register(user *model.User) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}
	user.Password = string(hashed)
	return u.repo.CreateUser(user)
}

func (u *userUsecase) Login(email, password string) (model.User, error) {
	user, err := u.repo.GetUserByEmail(email)
	if err != nil {
		return model.User{}, err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return model.User{}, bcrypt.ErrMismatchedHashAndPassword
	}
	return user, nil
}

func (u *userUsecase) GetProfile(id int) (model.User, error) {
	return u.repo.GetUserByID(id)
}

func (u *userUsecase) UpdateUser(user model.User) error {
	return u.repo.UpdateUser(user)
}

func (u *userUsecase) DeleteUser(id int) error {
	return u.repo.DeleteUser(id)
}
