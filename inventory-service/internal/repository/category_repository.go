package repository

import (
	"database/sql"
	"inventory/internal/model"
)

type CategoryRepository interface {
	GetAll() ([]model.Category, error)
	Create(category model.Category) error
	Update(id int, category model.Category) error
	Delete(id int) error
}

type categoryRepo struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepo{db}
}

func (r *categoryRepo) GetAll() ([]model.Category, error) {
	rows, err := r.db.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func (r *categoryRepo) Create(category model.Category) error {
	_, err := r.db.Exec("INSERT INTO categories (name) VALUES ($1)", category.Name)
	return err
}

func (r *categoryRepo) Update(id int, category model.Category) error {
	_, err := r.db.Exec("UPDATE categories SET name=$1 WHERE id=$2", category.Name, id)
	return err
}

func (r *categoryRepo) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM categories WHERE id=$1", id)
	return err
}
