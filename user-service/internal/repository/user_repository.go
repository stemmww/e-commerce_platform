package repository

import "user-service/internal/model"

type UserRepository interface {
	Register(user *model.User) error
	Authenticate(email, password string) (*model.User, error)
	GetByID(id string) (*model.User, error)
}
