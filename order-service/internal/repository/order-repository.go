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
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert order
	query := `INSERT INTO orders (user_id, total, status) VALUES ($1, $2, $3) RETURNING id`
	err = tx.QueryRow(query, order.UserID, order.Total, order.Status).Scan(&order.ID)
	if err != nil {
		return err
	}

	// Insert order items
	for _, item := range order.OrderItems {
		itemQuery := `INSERT INTO order_items (order_id, product_id, quantity) VALUES ($1, $2, $3)`
		_, err := tx.Exec(itemQuery, order.ID, item.ProductID, item.Quantity)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *orderRepository) GetAll() ([]model.Order, error) {
	rows, err := r.db.Query(`SELECT id, user_id, total, status FROM orders`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Total, &o.Status); err != nil {
			return nil, err
		}

		// Fetch order items for this order
		itemRows, err := r.db.Query(`SELECT id, order_id, product_id, quantity FROM order_items WHERE order_id = $1`, o.ID)
		if err != nil {
			return nil, err
		}
		defer itemRows.Close()

		var items []model.OrderItem
		for itemRows.Next() {
			var item model.OrderItem
			if err := itemRows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity); err != nil {
				return nil, err
			}
			items = append(items, item)
		}

		o.OrderItems = items
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *orderRepository) GetByID(id int) (*model.Order, error) {
	row := r.db.QueryRow(`SELECT id, user_id, total, status FROM orders WHERE id = $1`, id)
	var o model.Order
	if err := row.Scan(&o.ID, &o.UserID, &o.Total, &o.Status); err != nil {
		return nil, err
	}

	itemRows, err := r.db.Query(`SELECT id, order_id, product_id, quantity FROM order_items WHERE order_id = $1`, o.ID)
	if err != nil {
		return nil, err
	}
	defer itemRows.Close()

	var items []model.OrderItem
	for itemRows.Next() {
		var item model.OrderItem
		if err := itemRows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	o.OrderItems = items
	return &o, nil
}

func (r *orderRepository) UpdateStatus(id int, status string) error {
	_, err := r.db.Exec(`UPDATE orders SET status = $1 WHERE id = $2`, status, id)
	return err
}
