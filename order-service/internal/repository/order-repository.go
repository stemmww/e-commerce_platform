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
	Delete(id int) error
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

	_, err = tx.Exec(`INSERT INTO orders (id, user_id, status) VALUES ($1, $2, $3)`,
		order.ID, order.UserID, order.Status)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, item := range order.Items {
		_, err := tx.Exec(`INSERT INTO order_items (order_id, product_id, quantity) VALUES ($1, $2, $3)`,
			order.ID, item.ProductID, item.Quantity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *orderRepository) GetAll() ([]model.Order, error) {
	rows, err := r.db.Query(`SELECT id, user_id, status FROM orders`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Status); err != nil {
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

		o.Items = items
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *orderRepository) GetByID(id int) (*model.Order, error) {
	row := r.db.QueryRow(`SELECT id, user_id, status FROM orders WHERE id = $1`, id)
	var o model.Order
	if err := row.Scan(&o.ID, &o.UserID, &o.Status); err != nil {
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

	o.Items = items
	return &o, nil
}

func (r *orderRepository) UpdateStatus(id int, status string) error {
	_, err := r.db.Exec(`UPDATE orders SET status = $1 WHERE id = $2`, status, id)
	return err
}

func (r *orderRepository) Delete(id int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// First delete items
	_, err = tx.Exec(`DELETE FROM order_items WHERE order_id = $1`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Then delete order
	_, err = tx.Exec(`DELETE FROM orders WHERE id = $1`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
