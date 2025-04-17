package repository

import (
	"database/sql"
	"order-service/internal/model"
)

type OrderRepository interface {
	Create(order *model.Order) error
	GetAll() ([]model.Order, error)
	GetByID(id int) (*model.Order, error)
	UpdateStatus(id int, status string) error
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *model.Order) error {
	query := `INSERT INTO orders (user_id, product_id, quantity, status) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, order.UserID, order.ProductID, order.Quantity, order.Status)
	return err
}

func (r *orderRepository) GetAll() ([]model.Order, error) {
	rows, err := r.db.Query(`SELECT id, user_id, product_id, quantity, status FROM orders`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.ProductID, &o.Quantity, &o.Status); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *orderRepository) GetByID(id int) (*model.Order, error) {
	row := r.db.QueryRow(`SELECT id, user_id, product_id, quantity, status FROM orders WHERE id = $1`, id)
	var o model.Order
	err := row.Scan(&o.ID, &o.UserID, &o.ProductID, &o.Quantity, &o.Status)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *orderRepository) UpdateStatus(id int, status string) error {
	_, err := r.db.Exec(`UPDATE orders SET status = $1 WHERE id = $2`, status, id)
	return err
}
