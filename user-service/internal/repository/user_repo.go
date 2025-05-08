package repository

import (
	"user-service/internal/model"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByEmail(email string) (model.User, error)
	GetUserByID(id int) (model.User, error)
	UpdateUser(user model.User) error
	DeleteUser(id int) error
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(user *model.User) error {
	return r.db.QueryRowx(
		`INSERT INTO users (email, password, name) VALUES ($1, $2, $3) RETURNING id`,
		user.Email, user.Password, user.Name,
	).Scan(&user.ID)
}

func (r *userRepo) GetUserByEmail(email string) (model.User, error) {
	var user model.User
	err := r.db.Get(&user, `SELECT * FROM users WHERE email=$1`, email)
	return user, err
}

func (r *userRepo) GetUserByID(id int) (model.User, error) {
	var user model.User
	err := r.db.Get(&user, `SELECT id, email, name FROM users WHERE id=$1`, id)
	return user, err
}

func (r *userRepo) UpdateUser(user model.User) error {
	_, err := r.db.Exec(`UPDATE users SET name=$1, email=$2 WHERE id=$3`, user.Name, user.Email, user.ID)
	return err
}

func (r *userRepo) DeleteUser(id int) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id=$1`, id)
	return err
}
