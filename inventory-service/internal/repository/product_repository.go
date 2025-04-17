package repository

import (
	"database/sql"
	"inventory/internal/model"
)

type ProductRepository interface {
	Create(product model.Product) error
	GetByID(id int) (*model.Product, error)
	GetAll() ([]model.Product, error)
	Update(product model.Product) error
	Delete(id int) error
}

type productRepo struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepo{db}
}

func (r *productRepo) Create(product model.Product) error {
	_, err := r.db.Exec(
		"INSERT INTO products (name, category_id, price, stock, description) VALUES ($1, $2, $3, $4, $5)",
		product.Name, product.CategoryID, product.Price, product.Stock, product.Description)
	return err
}

func (r *productRepo) GetByID(id int) (*model.Product, error) {
	row := r.db.QueryRow("SELECT id, name, category_id, price, stock, description FROM products WHERE id = $1", id)
	var p model.Product
	err := row.Scan(&p.ID, &p.Name, &p.CategoryID, &p.Price, &p.Stock, &p.Description)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *productRepo) GetAll() ([]model.Product, error) {
	rows, err := r.db.Query("SELECT id, name, category_id, price, stock, description FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		err := rows.Scan(&p.ID, &p.Name, &p.CategoryID, &p.Price, &p.Stock, &p.Description)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *productRepo) Update(product model.Product) error {
	_, err := r.db.Exec(
		"UPDATE products SET name = $1, category_id = $2, price = $3, stock = $4, description = $5 WHERE id = $6",
		product.Name, product.CategoryID, product.Price, product.Stock, product.Description, product.ID)
	return err
}

func (r *productRepo) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM products WHERE id = $1", id)
	return err
}
