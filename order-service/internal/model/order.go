package model

type Order struct {
	ID         int         `json:"id"`
	UserID     int         `json:"user_id"`
	Total      float64     `json:"total"`
	Status     string      `json:"status"`
	OrderItems []OrderItem `json:"order_items"` // slice of related items
}

type OrderItem struct {
	ID        int `json:"id"`
	OrderID   int `json:"order_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
